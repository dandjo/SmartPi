package config

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nDenerserve/SmartPi/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Moduleconfig struct {
	// [base]
	I2CDevice string
	Webserver bool
	Vfs       bool
	LogLevel  log.Level

	// [webserver]
	WebserverPort int
	AppKey        string

	// [digitalout]
	AllowedDigitalOutUser []string

	// [etemperature]
	AllowedEtemperatureUser       []string
	EtemperatureI2CAddress        uint16
	EtemperatureSamplerate        int
	EtemperatureSharedFileEnabled bool
	EtemperatureSharedDir         string
	EtemperatureSharedFile        string

	// [lorawan]
	LoRaWANEnabled             bool
	LoRaWANSharedDirs          []string
	LoRaWANSharedFilesElements []string
	LoRaWANSerialPort          string
	LoRaWANSendInterval        int
	LoRaWANApplicationEUI      string
	LoRaWANApplicationKey      string
	LoRaWANDataRate            int

	// [LTE]
	LteEnabled bool
	LteAPN     string
	LteSimPin  string

	// s := strings.Split("a,b,c", ",")
}

var mcfg *ini.File
var merr error

func (p *Moduleconfig) ReadParameterFromFile() {

	mcfg, merr = ini.Load("/etc/smartpi_modules")
	if merr != nil {
		panic(merr)
	}

	// [base]
	p.I2CDevice = mcfg.Section("base").Key("i2c_device").MustString("/dev/i2c-1")
	p.Webserver = mcfg.Section("base").Key("webserver").MustBool(true)
	p.Vfs = mcfg.Section("base").Key("vfs").MustBool(true)
	p.LogLevel, merr = log.ParseLevel(mcfg.Section("base").Key("loglevel").MustString("info"))
	if merr != nil {
		panic(merr)
	}

	// [webserver]
	p.WebserverPort = mcfg.Section("webserver").Key("port").MustInt(1080)
	p.AppKey = mcfg.Section("webserver").Key("appkey").MustString("ew980723j35h546ergr97fqw4!234490#t33465")

	// [digitalout]
	p.AllowedDigitalOutUser = strings.Split(mcfg.Section("digitalout").Key("allowed_user").String(), ",")

	// [etemperature]
	p.AllowedEtemperatureUser = strings.Split(mcfg.Section("etemperature").Key("allowed_user").String(), ",")
	if p.EtemperatureI2CAddress, err = utils.DecodeUint16(mcfg.Section("etemperature").Key("i2c_address").MustString("0x52")); err != nil {
		log.Fatal(err)
	}
	p.EtemperatureSamplerate = mcfg.Section("etemperature").Key("samplerate").MustInt(6)
	p.EtemperatureSharedFileEnabled = mcfg.Section("etemperature").Key("shared_file_enabled").MustBool(true)
	p.EtemperatureSharedDir = mcfg.Section("etemperature").Key("shared_dir").MustString("/var/run")
	p.EtemperatureSharedFile = mcfg.Section("etemperature").Key("shared_file").MustString("smartpi_etemperature_values")

	// [lorawan]
	p.LoRaWANEnabled = mcfg.Section("lorawan").Key("lorawan_enabled").MustBool(true)
	p.LoRaWANSharedDirs = strings.Split(mcfg.Section("lorawan").Key("shared_files_path").String(), ",")
	if len(p.LoRaWANSharedDirs) == 0 {
		p.LoRaWANSharedDirs = append(p.LoRaWANSharedDirs, "/var/run/smartpi_values")
	}
	p.LoRaWANSharedFilesElements = strings.Split(mcfg.Section("lorawan").Key("shared_files_elements").String(), ",")
	if len(p.LoRaWANSharedFilesElements) == 0 {
		p.LoRaWANSharedFilesElements = append(p.LoRaWANSharedFilesElements, "1:2f:1.0|2:2f:1.0|3:2f:1.0|4:2f:1.0|5:2f:1.0|6:2f:1.0|7:2f:1.0,1:2f:1.0|2:2f:1.0|3:2f:1.0|4:2f:1.0|5:2f:1.0|6:2f:1.0|7:2f:1.0")
	}
	p.LoRaWANSendInterval = mcfg.Section("lorawan").Key("interval").MustInt(60)
	p.LoRaWANSerialPort = mcfg.Section("lorawan").Key("serial_port").MustString("/dev/ttyS0")
	p.LoRaWANApplicationEUI = mcfg.Section("lorawan").Key("applicationEUI").MustString("")
	p.LoRaWANApplicationKey = mcfg.Section("lorawan").Key("applicationKey").MustString("")
	p.LoRaWANDataRate = mcfg.Section("lorawan").Key("datarate").MustInt(5)

	// [LTE]
	p.LteEnabled = mcfg.Section("lte").Key("enabled").MustBool(true)

}

func (p *Moduleconfig) SaveParameterToFile() {

	// [base]
	_, merr = mcfg.Section("base").NewKey("i2c_device", p.I2CDevice)
	_, merr = mcfg.Section("base").NewKey("webserver", strconv.FormatBool(p.Webserver))
	_, merr = mcfg.Section("base").NewKey("vfs", strconv.FormatBool(p.Vfs))
	_, merr = mcfg.Section("base").NewKey("loglevel", p.LogLevel.String())

	// [webserver]
	_, merr = mcfg.Section("webserver").NewKey("port", strconv.FormatInt(int64(p.WebserverPort), 10))
	_, merr = mcfg.Section("appkey").NewKey("appkey", p.AppKey)

	// [digitalout]
	_, merr = mcfg.Section("digitalout").NewKey("allowed_user", strings.Join(p.AllowedDigitalOutUser, ","))

	//[etemperature]
	_, merr = mcfg.Section("etemperature").NewKey("allowed_user", strings.Join(p.AllowedEtemperatureUser, ","))
	_, merr = mcfg.Section("etemperature").NewKey("i2c_address", utils.EncodeUint64(uint64(p.EtemperatureI2CAddress)))
	_, merr = mcfg.Section("etemperature").NewKey("samplerate", strconv.FormatInt(int64(p.EtemperatureSamplerate), 10))
	_, merr = mcfg.Section("etemperature").NewKey("shared_file_enabled", strconv.FormatBool(p.EtemperatureSharedFileEnabled))
	_, merr = mcfg.Section("etemperature").NewKey("shared_dir", p.EtemperatureSharedDir)
	_, merr = mcfg.Section("etemperature").NewKey("shared_file", p.EtemperatureSharedFile)

	//[lorawan]
	_, merr = mcfg.Section("lorawan").NewKey("shared_file_enabled", strconv.FormatBool(p.LoRaWANEnabled))
	_, merr = mcfg.Section("lorawan").NewKey("shared_files_path", strings.Join(p.LoRaWANSharedDirs[:], ","))
	_, merr = mcfg.Section("lorawan").NewKey("shared_files_elements", strings.Join(p.LoRaWANSharedFilesElements[:], ","))
	_, merr = mcfg.Section("lorawan").NewKey("interval", strconv.FormatInt(int64(p.EtemperatureSamplerate), 60))
	_, merr = mcfg.Section("lorawan").NewKey("serial_port", p.LoRaWANSerialPort)
	_, merr = mcfg.Section("lorawan").NewKey("applicationEUI", p.LoRaWANApplicationEUI)
	_, merr = mcfg.Section("lorawan").NewKey("applicationKey", p.LoRaWANApplicationKey)
	_, merr = mcfg.Section("lorawan").NewKey("datarate", strconv.FormatInt(int64(p.LoRaWANDataRate), 10))

	tmpFile := "/tmp/smartpi_modules"
	merr := mcfg.SaveTo(tmpFile)
	if merr != nil {
		panic(merr)
	}

	srcFile, merr := os.Open(tmpFile)
	utils.Checklog(merr)
	defer srcFile.Close()

	destFile, merr := os.Create("/etc/smartpi_modules") // creates if file doesn't exist
	utils.Checklog(merr)
	defer destFile.Close()

	_, merr = io.Copy(destFile, srcFile)
	utils.Checklog(merr)

	merr = destFile.Sync()
	utils.Checklog(merr)

	defer os.Remove(tmpFile)
}

func NewModuleconfig() *Moduleconfig {

	t := new(Moduleconfig)
	t.ReadParameterFromFile()
	return t
}
