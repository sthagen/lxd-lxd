//go:build linux && cgo && !agent

package sys

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/syndtr/gocapability/capability"

	"github.com/canonical/lxd/lxd/db/cluster"
	"github.com/canonical/lxd/lxd/db/warningtype"
	"github.com/canonical/lxd/lxd/util"
	"github.com/canonical/lxd/shared"
	"github.com/canonical/lxd/shared/logger"
)

// Initialize AppArmor-specific attributes.
func (s *OS) initAppArmor() []cluster.Warning {
	var dbWarnings []cluster.Warning

	/* Detect AppArmor availability */
	_, err := exec.LookPath("apparmor_parser")
	if shared.IsFalse(os.Getenv("LXD_SECURITY_APPARMOR")) {
		logger.Warnf("AppArmor support has been manually disabled")
		dbWarnings = append(dbWarnings, cluster.Warning{
			TypeCode:    warningtype.AppArmorNotAvailable,
			LastMessage: "Manually disabled",
		})
	} else if !shared.IsDir("/sys/kernel/security/apparmor") {
		logger.Warnf("AppArmor support has been disabled because of lack of kernel support")
		dbWarnings = append(dbWarnings, cluster.Warning{
			TypeCode:    warningtype.AppArmorNotAvailable,
			LastMessage: "Disabled because of lack of kernel support",
		})
	} else if err != nil {
		logger.Warnf("AppArmor support has been disabled because 'apparmor_parser' couldn't be found")
		dbWarnings = append(dbWarnings, cluster.Warning{
			TypeCode:    warningtype.AppArmorNotAvailable,
			LastMessage: "Disabled because 'apparmor_parser' couldn't be found",
		})
	} else {
		s.AppArmorAvailable = true
	}

	/* Detect AppArmor stacking support */
	s.AppArmorStacking = appArmorCanStack()

	/* Detect existing AppArmor stack */
	if shared.PathExists("/sys/kernel/security/apparmor/.ns_stacked") {
		contentBytes, err := os.ReadFile("/sys/kernel/security/apparmor/.ns_stacked")
		if err == nil && string(contentBytes) == "yes\n" {
			s.AppArmorStacked = true
		}
	}

	/* Detect AppArmor admin support */
	if !haveMacAdmin() {
		if s.AppArmorAvailable {
			logger.Warnf("Per-container AppArmor profiles are disabled because the mac_admin capability is missing")
		}
	} else if s.RunningInUserNS && !s.AppArmorStacked {
		if s.AppArmorAvailable {
			logger.Warnf("Per-container AppArmor profiles are disabled because LXD is running in an unprivileged container without stacking")
		}
	} else {
		s.AppArmorAdmin = true
	}

	/* Detect AppArmor confinment */
	profile := util.AppArmorProfile()
	/* if AppArmor is enabled on the system but there is no profile then
	   the application has the profile name "unconfined", or if AppArmor
	   is not supported then the label will be empty and finally newer
	   AppArmor releases support a profile mode "unconfined" where an
	   application can have a profile defined for it which effectively
	   is equivalent to the traditional "unconfined" label - in this case
	   the label will have the suffix " (unconfined)" */
	if profile != "unconfined" && profile != "" && !strings.HasSuffix(profile, "(unconfined)") {
		if s.AppArmorAvailable {
			logger.Warnf("Per-container AppArmor profiles are disabled because LXD is already protected by AppArmor via profile %q", profile)
		}

		s.AppArmorConfined = true
	}

	s.AppArmorFeatures.Map = map[string]bool{}

	return dbWarnings
}

func haveMacAdmin() bool {
	c, err := capability.NewPid2(0)
	if err != nil {
		return false
	}

	err = c.Load()
	if err != nil {
		return false
	}

	if c.Get(capability.EFFECTIVE, capability.CAP_MAC_ADMIN) {
		return true
	}

	return false
}

// Returns true if AppArmor stacking support is available.
func appArmorCanStack() bool {
	contentBytes, err := os.ReadFile("/sys/kernel/security/apparmor/features/domain/stack")
	if err != nil {
		return false
	}

	if string(contentBytes) != "yes\n" {
		return false
	}

	contentBytes, err = os.ReadFile("/sys/kernel/security/apparmor/features/domain/version")
	if err != nil {
		return false
	}

	content := string(contentBytes)

	parts := strings.Split(strings.TrimSpace(content), ".")

	if len(parts) == 0 {
		logger.Warn("Unknown apparmor domain version", logger.Ctx{"version": content})
		return false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.Warn("Unknown apparmor domain version", logger.Ctx{"version": content})
		return false
	}

	minor := 0
	if len(parts) == 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			logger.Warn("Unknown apparmor domain version", logger.Ctx{"version": content})
			return false
		}
	}

	return major >= 1 && minor >= 2
}
