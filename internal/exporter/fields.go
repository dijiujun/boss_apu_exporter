package exporter

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	/*uuidQField QField = "uuid"*/
	cardQField        QField = "card_index"
	nameQField        QField = "name"
	softVersionQField QField = "soft_version"
	apuVersionQField  QField = "apu_version"
	boardIDQField     QField = "board_id"
	chipTypeQField    QField = "chip_type"
	busIDQField       QField = "bus_id"
	die0SNQField      QField = "die0_sn"
	die1SNQField      QField = "die1_sn"
	die2SNQField      QField = "die2_sn"
	die3SNQField      QField = "die3_sn"

	/*driverModelCurrentQField QField = "driver_model.current"
	driverModelPendingQField QField = "driver_model.pending"
	vBiosVersionQField       QField = "vbios_version"
	driverVersionQField      QField = "driver_version"*/
	qFieldsAuto   = "AUTO"
	DefaultQField = qFieldsAuto
)

var (
	ErrNoQueryFields = errors.New("could not extract any query fields")

	fieldRegex = regexp.MustCompile(`(?m)\n\s*\n^"([^"]+)"`)

	//nolint:gochecknoglobals
	fallbackQFieldToRFieldMap = map[QField]RField{
		"timestamp":           "timestamp",
		"card_index":          "card_index",
		"bus_id":              "bus_id",
		"link_width":          "link_width",
		"soft_version":        "soft_version",
		"board_id":            "board_id",
		"cpld_version":        "cpld_version",
		"name":                "name",
		"chip_type":           "chip_type",
		"apu_version":         "apu_version",
		"pwr":                 "pwr",
		"inlet_temp":          "inlet_temp",
		"outlet_temp":         "outlet_temp",
		"card_temp":           "card_temp",
		"die0_sn":             "die0_sn",
		"die0_cpu_usage_rate": "die0_cpu_usage_rate",
		"die0_gpu_usage_rate": "die0_gpu_usage_rate",
		"die0_bxlink_temp":    "die0_bxlink_temp",
		"die0_gpgpu_temp":     "die0_gpgpu_temp",
		"die0_cpu_temp":       "die0_cpu_temp",
		"die0_ddr_temp":       "die0_ddr_temp",
		"die0_decoder_temp":   "die0_decoder_temp",
		"die0_sys_temp":       "die0_sys_temp",
		"die0_slowdown_temp":  "die0_slowdown_temp",
		"die0_shutdown_temp":  "die0_shutdown_temp",
		"die0_gpu_clock":      "die0_gpu_clock",
		"die0_cpu_clock":      "die0_cpu_clock",
		"die0_ddr_clock":      "die0_ddr_clock",
		"die0_mem_total":      "die0_mem_total",
		"die0_mem_used":       "die0_mem_used",
		"die0_mem_free":       "die0_mem_free",
		"die0_gpu_mem_total":  "die0_gpu_mem_total",
		"die0_gpu_mem_used":   "die0_gpu_mem_used",
		"die0_gpu_mem_free":   "die0_gpu_mem_free",
		"die0_phy_mem":        "die0_phy_mem",
		"die1_sn":             "die1_sn",
		"die1_cpu_usage_rate": "die1_cpu_usage_rate",
		"die1_gpu_usage_rate": "die1_gpu_usage_rate",
		"die1_bxlink_temp":    "die1_bxlink_temp",
		"die1_gpgpu_temp":     "die1_gpgpu_temp",
		"die1_cpu_temp":       "die1_cpu_temp",
		"die1_ddr_temp":       "die1_ddr_temp",
		"die1_decoder_temp":   "die1_decoder_temp",
		"die1_sys_temp":       "die1_sys_temp",
		"die1_slowdown_temp":  "die1_slowdown_temp",
		"die1_shutdown_temp":  "die1_shutdown_temp",
		"die1_gpu_clock":      "die1_gpu_clock",
		"die1_cpu_clock":      "die1_cpu_clock",
		"die1_ddr_clock":      "die1_ddr_clock",
		"die1_mem_total":      "die1_mem_total",
		"die1_mem_used":       "die1_mem_used",
		"die1_mem_free":       "die1_mem_free",
		"die1_gpu_mem_total":  "die1_gpu_mem_total",
		"die1_gpu_mem_used":   "die1_gpu_mem_used",
		"die1_gpu_mem_free":   "die1_gpu_mem_free",
		"die1_phy_mem":        "die1_phy_mem",
		"die2_sn":             "die2_sn",
		"die2_cpu_usage_rate": "die2_cpu_usage_rate",
		"die2_gpu_usage_rate": "die2_gpu_usage_rate",
		"die2_bxlink_temp":    "die2_bxlink_temp",
		"die2_gpgpu_temp":     "die2_gpgpu_temp",
		"die2_cpu_temp":       "die2_cpu_temp",
		"die2_ddr_temp":       "die2_ddr_temp",
		"die2_decoder_temp":   "die2_decoder_temp",
		"die2_sys_temp":       "die2_sys_temp",
		"die2_slowdown_temp":  "die2_slowdown_temp",
		"die2_shutdown_temp":  "die2_shutdown_temp",
		"die2_gpu_clock":      "die2_gpu_clock",
		"die2_cpu_clock":      "die2_cpu_clock",
		"die2_ddr_clock":      "die2_ddr_clock",
		"die2_mem_total":      "die2_mem_total",
		"die2_mem_used":       "die2_mem_used",
		"die2_mem_free":       "die2_mem_free",
		"die2_gpu_mem_total":  "die2_gpu_mem_total",
		"die2_gpu_mem_used":   "die2_gpu_mem_used",
		"die2_gpu_mem_free":   "die2_gpu_mem_free",
		"die2_phy_mem":        "die2_phy_mem",
		"die3_sn":             "die3_sn",
		"die3_cpu_usage_rate": "die3_cpu_usage_rate",
		"die3_gpu_usage_rate": "die3_gpu_usage_rate",
		"die3_bxlink_temp":    "die3_bxlink_temp",
		"die3_gpgpu_temp":     "die3_gpgpu_temp",
		"die3_cpu_temp":       "die3_cpu_temp",
		"die3_ddr_temp":       "die3_ddr_temp",
		"die3_decoder_temp":   "die3_decoder_temp",
		"die3_sys_temp":       "die3_sys_temp",
		"die3_slowdown_temp":  "die3_slowdown_temp",
		"die3_shutdown_temp":  "die3_shutdown_temp",
		"die3_gpu_clock":      "die3_gpu_clock",
		"die3_cpu_clock":      "die3_cpu_clock",
		"die3_ddr_clock":      "die3_ddr_clock",
		"die3_mem_total":      "die3_mem_total",
		"die3_mem_used":       "die3_mem_used",
		"die3_mem_free":       "die3_mem_free",
		"die3_gpu_mem_total":  "die3_gpu_mem_total",
		"die3_gpu_mem_used":   "die3_gpu_mem_used",
		"die3_gpu_mem_free":   "die3_gpu_mem_free",
		"die3_phy_mem":        "die3_phy_mem",
	}
	/*fallbackQFieldToRFieldMap = map[QField]RField{
		"timestamp":                         "timestamp",
		"driver_version":                    "driver_version",
		"count":                             "count",
		"name":                              "name",
		"serial":                            "serial",
		"uuid":                              "uuid",
		"pci.bus_id":                        "pci.bus_id",
		"pci.domain":                        "pci.domain",
		"pci.bus":                           "pci.bus",
		"pci.device":                        "pci.device",
		"pci.device_id":                     "pci.device_id",
		"pci.sub_device_id":                 "pci.sub_device_id",
		"pcie.link.gen.current":             "pcie.link.gen.current",
		"pcie.link.gen.max":                 "pcie.link.gen.max",
		"pcie.link.width.current":           "pcie.link.width.current",
		"pcie.link.width.max":               "pcie.link.width.max",
		"index":                             "index",
		"display_mode":                      "display_mode",
		"display_active":                    "display_active",
		"persistence_mode":                  "persistence_mode",
		"accounting.mode":                   "accounting.mode",
		"accounting.buffer_size":            "accounting.buffer_size",
		"driver_model.current":              "driver_model.current",
		"driver_model.pending":              "driver_model.pending",
		"vbios_version":                     "vbios_version",
		"inforom.img":                       "inforom.img",
		"inforom.oem":                       "inforom.oem",
		"inforom.ecc":                       "inforom.ecc",
		"inforom.pwr":                       "inforom.pwr",
		"gom.current":                       "gom.current",
		"gom.pending":                       "gom.pending",
		"fan.speed":                         "fan.speed [%]",
		"pstate":                            "pstate",
		"clocks_throttle_reasons.supported": "clocks_throttle_reasons.supported",
		"clocks_throttle_reasons.active":    "clocks_throttle_reasons.active",
		"clocks_throttle_reasons.gpu_idle":  "clocks_throttle_reasons.gpu_idle",
		"clocks_throttle_reasons.applications_clocks_setting": "clocks_throttle_reasons.applications_clocks_setting",
		"clocks_throttle_reasons.sw_power_cap":                "clocks_throttle_reasons.sw_power_cap",
		"clocks_throttle_reasons.hw_slowdown":                 "clocks_throttle_reasons.hw_slowdown",
		"clocks_throttle_reasons.hw_thermal_slowdown":         "clocks_throttle_reasons.hw_thermal_slowdown",
		"clocks_throttle_reasons.hw_power_brake_slowdown":     "clocks_throttle_reasons.hw_power_brake_slowdown",
		"clocks_throttle_reasons.sw_thermal_slowdown":         "clocks_throttle_reasons.sw_thermal_slowdown",
		"clocks_throttle_reasons.sync_boost":                  "clocks_throttle_reasons.sync_boost",
		"memory.total":                                        "memory.total [MiB]",
		"memory.used":                                         "memory.used [MiB]",
		"memory.free":                                         "memory.free [MiB]",
		"compute_mode":                                        "compute_mode",
		"utilization.gpu":                                     "utilization.gpu [%]",
		"utilization.memory":                                  "utilization.memory [%]",
		"encoder.stats.sessionCount":                          "encoder.stats.sessionCount",
		"encoder.stats.averageFps":                            "encoder.stats.averageFps",
		"encoder.stats.averageLatency":                        "encoder.stats.averageLatency",
		"ecc.mode.current":                                    "ecc.mode.current",
		"ecc.mode.pending":                                    "ecc.mode.pending",
		"ecc.errors.corrected.volatile.device_memory":         "ecc.errors.corrected.volatile.device_memory",
		"ecc.errors.corrected.volatile.dram":                  "ecc.errors.corrected.volatile.dram",
		"ecc.errors.corrected.volatile.register_file":         "ecc.errors.corrected.volatile.register_file",
		"ecc.errors.corrected.volatile.l1_cache":              "ecc.errors.corrected.volatile.l1_cache",
		"ecc.errors.corrected.volatile.l2_cache":              "ecc.errors.corrected.volatile.l2_cache",
		"ecc.errors.corrected.volatile.texture_memory":        "ecc.errors.corrected.volatile.texture_memory",
		"ecc.errors.corrected.volatile.cbu":                   "ecc.errors.corrected.volatile.cbu",
		"ecc.errors.corrected.volatile.sram":                  "ecc.errors.corrected.volatile.sram",
		"ecc.errors.corrected.volatile.total":                 "ecc.errors.corrected.volatile.total",
		"ecc.errors.corrected.aggregate.device_memory":        "ecc.errors.corrected.aggregate.device_memory",
		"ecc.errors.corrected.aggregate.dram":                 "ecc.errors.corrected.aggregate.dram",
		"ecc.errors.corrected.aggregate.register_file":        "ecc.errors.corrected.aggregate.register_file",
		"ecc.errors.corrected.aggregate.l1_cache":             "ecc.errors.corrected.aggregate.l1_cache",
		"ecc.errors.corrected.aggregate.l2_cache":             "ecc.errors.corrected.aggregate.l2_cache",
		"ecc.errors.corrected.aggregate.texture_memory":       "ecc.errors.corrected.aggregate.texture_memory",
		"ecc.errors.corrected.aggregate.cbu":                  "ecc.errors.corrected.aggregate.cbu",
		"ecc.errors.corrected.aggregate.sram":                 "ecc.errors.corrected.aggregate.sram",
		"ecc.errors.corrected.aggregate.total":                "ecc.errors.corrected.aggregate.total",
		"ecc.errors.uncorrected.volatile.device_memory":       "ecc.errors.uncorrected.volatile.device_memory",
		"ecc.errors.uncorrected.volatile.dram":                "ecc.errors.uncorrected.volatile.dram",
		"ecc.errors.uncorrected.volatile.register_file":       "ecc.errors.uncorrected.volatile.register_file",
		"ecc.errors.uncorrected.volatile.l1_cache":            "ecc.errors.uncorrected.volatile.l1_cache",
		"ecc.errors.uncorrected.volatile.l2_cache":            "ecc.errors.uncorrected.volatile.l2_cache",
		"ecc.errors.uncorrected.volatile.texture_memory":      "ecc.errors.uncorrected.volatile.texture_memory",
		"ecc.errors.uncorrected.volatile.cbu":                 "ecc.errors.uncorrected.volatile.cbu",
		"ecc.errors.uncorrected.volatile.sram":                "ecc.errors.uncorrected.volatile.sram",
		"ecc.errors.uncorrected.volatile.total":               "ecc.errors.uncorrected.volatile.total",
		"ecc.errors.uncorrected.aggregate.device_memory":      "ecc.errors.uncorrected.aggregate.device_memory",
		"ecc.errors.uncorrected.aggregate.dram":               "ecc.errors.uncorrected.aggregate.dram",
		"ecc.errors.uncorrected.aggregate.register_file":      "ecc.errors.uncorrected.aggregate.register_file",
		"ecc.errors.uncorrected.aggregate.l1_cache":           "ecc.errors.uncorrected.aggregate.l1_cache",
		"ecc.errors.uncorrected.aggregate.l2_cache":           "ecc.errors.uncorrected.aggregate.l2_cache",
		"ecc.errors.uncorrected.aggregate.texture_memory":     "ecc.errors.uncorrected.aggregate.texture_memory",
		"ecc.errors.uncorrected.aggregate.cbu":                "ecc.errors.uncorrected.aggregate.cbu",
		"ecc.errors.uncorrected.aggregate.sram":               "ecc.errors.uncorrected.aggregate.sram",
		"ecc.errors.uncorrected.aggregate.total":              "ecc.errors.uncorrected.aggregate.total",
		"retired_pages.single_bit_ecc.count":                  "retired_pages.single_bit_ecc.count",
		"retired_pages.double_bit.count":                      "retired_pages.double_bit.count",
		"retired_pages.pending":                               "retired_pages.pending",
		"temperature.gpu":                                     "temperature.gpu",
		"temperature.memory":                                  "temperature.memory",
		"power.management":                                    "power.management",
		"power.draw":                                          "power.draw [W]",
		"power.limit":                                         "power.limit [W]",
		"enforced.power.limit":                                "enforced.power.limit [W]",
		"power.default_limit":                                 "power.default_limit [W]",
		"power.min_limit":                                     "power.min_limit [W]",
		"power.max_limit":                                     "power.max_limit [W]",
		"clocks.current.graphics":                             "clocks.current.graphics [MHz]",
		"clocks.current.sm":                                   "clocks.current.sm [MHz]",
		"clocks.current.memory":                               "clocks.current.memory [MHz]",
		"clocks.current.video":                                "clocks.current.video [MHz]",
		"clocks.applications.graphics":                        "clocks.applications.graphics [MHz]",
		"clocks.applications.memory":                          "clocks.applications.memory [MHz]",
		"clocks.default_applications.graphics":                "clocks.default_applications.graphics [MHz]",
		"clocks.default_applications.memory":                  "clocks.default_applications.memory [MHz]",
		"clocks.max.graphics":                                 "clocks.max.graphics [MHz]",
		"clocks.max.sm":                                       "clocks.max.sm [MHz]",
		"clocks.max.memory":                                   "clocks.max.memory [MHz]",
		"mig.mode.current":                                    "mig.mode.current",
		"mig.mode.pending":                                    "mig.mode.pending",
	}*/
)

func ParseAutoQFields(bossSmiCommand string, command runCmd) ([]QField, error) {
	cmdAndArgs := strings.Fields(bossSmiCommand)
	cmdAndArgs = append(cmdAndArgs, "--help-query-gpu")
	cmd := exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...) //nolint:gosec

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := command(cmd)

	outStr := stdout.String()
	errStr := stderr.String()

	exitCode := -1

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		exitCode = exitError.ExitCode()
	}

	if err != nil {
		return nil, fmt.Errorf("%w: command failed. code: %d | command: %s | stdout: %s | stderr: %s", err,
			exitCode, strings.Join(cmdAndArgs, " "), outStr, errStr)
	}

	fields := ExtractQFields(outStr)
	if fields == nil {
		return nil, fmt.Errorf("%w: code: %d | command: %s | stdout: %s | stderr: %s", ErrNoQueryFields,
			exitCode, strings.Join(cmdAndArgs, " "), outStr, errStr)
	}

	return fields, nil
}

func ExtractQFields(text string) []QField {
	found := fieldRegex.FindAllStringSubmatch(text, -1)

	fields := make([]QField, len(found))
	for i, ss := range found {
		fields[i] = QField(ss[1])
	}

	return fields
}

func toQFieldSlice(ss []string) []QField {
	r := make([]QField, len(ss))
	for i, s := range ss {
		r[i] = QField(s)
	}

	return r
}

func toRFieldSlice(ss []string) []RField {
	r := make([]RField, len(ss))
	for i, s := range ss {
		r[i] = RField(s)
	}

	return r
}

func QFieldSliceToStringSlice(qs []QField) []string {
	r := make([]string, len(qs))
	for i, q := range qs {
		r[i] = string(q)
	}

	return r
}
