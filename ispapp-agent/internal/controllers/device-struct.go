package controllers

import (
	"encoding/json"
)

type Device struct {
	Type                   string              `json:"type"`
	Login                  string              `json:"login"`
	Key                    string              `json:"key"`
	Asn                    string              `json:"asn"`
	Owner                  string              `json:"owner"`
	SimpleRotatedKey       string              `json:"simpleRotatedKey"`
	Uptime                 int64               `json:"uptime"`
	WanIp                  string              `json:"wanIp"`
	FwStatus               string              `json:"fwStatus"`
	ClientInfo             string              `json:"clientInfo"`
	Os                     string              `json:"os"`
	OsVersion              string              `json:"osVersion"`
	Fw                     string              `json:"fw,omitempty"`
	HardwareMake           string              `json:"hardwareMake"`
	HardwareModel          string              `json:"hardwareModel"`
	HardwareModelNumber    string              `json:"hardwareModelNumber"`
	HardwareSerialNumber   string              `json:"hardwareSerialNumber"`
	HardwareCpuInfo        string              `json:"hardwareCpuInfo"`
	OsBuildDate            string              `json:"osBuildDate"`
	Reboot                 int64               `json:"reboot"`
	Cmds                   []Cmd               `json:"cmds"`
	Cmd                    string              `json:"cmd"`
	Ws_Id                  string              `json:"ws_id"`
	UuidV4                 string              `json:"uuidv4"`
	Stdout                 string              `json:"stdout"`
	Stderr                 string              `json:"stderr"`
	Lat                    *float64            `json:"lat"`
	Lng                    *float64            `json:"lng"`
	Collectors             *Collector          `json:"collectors,omitempty"`
	Interfaces             []Interface         `json:"interfaces"`
	WirelessConfigured     []WirelessInterface `json:"wirelessConfigured"`
	WirelessConfig         string              `json:"wirelessConfig"`
	SecurityProfiles       []SecurityProfile   `json:"security-profiles"`
	WebshellSupport        bool                `json:"webshellSupport"`
	BandwidthTestSupport   bool                `json:"bandwidthTestSupport"`
	FirmwareUpgradeSupport bool                `json:"firmwareUpgradeSupport"`
	Hostname               string              `json:"hostname"`
	OutsideIP              string              `json:"outsideIP"`
	LastConfigRequest      int64               `json:"lastConfigRequest"`
	UsingWebSocket         bool                `json:"usingWebSocket"`
}

type Cmd struct {
	Cmd    string `json:"cmd"`
	Type   string `json:"type"`
	Ws_Id  string `json:"ws_id"`
	UuidV4 string `json:"uuidv4"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type SecurityProfile struct {
	ID                  *string  `json:".id,omitempty"`
	AuthenticationTypes []string `json:"authentication-types,omitempty"`
	Default             *bool    `json:"default,omitempty"`
	EAPMethods          []string `json:"eap-methods,omitempty"`
	GroupCiphers        []string `json:"group-ciphers,omitempty"`
	Mode                *string  `json:"mode,omitempty"`
	Name                *string  `json:"name,omitempty"`
	RadiusCalledFormat  *string  `json:"radius-called-format,omitempty"`
	Technology          *string  `json:"technology,omitempty"`
	WPAPreSharedKey     *string  `json:"wpa-pre-shared-key,omitempty"`
	WPA2PreSharedKey    *string  `json:"wpa2-pre-shared-key,omitempty"`
	WPA3PreSharedKey    *string  `json:"wpa3-pre-shared-key,omitempty"`
}

type Interface struct {
	If                     string   `json:"if"`
	DefaultIf              string   `json:"defaultIf"`
	Mac                    string   `json:"mac"`
	RecBytes               uint64   `json:"recBytes"`
	Up                     bool     `json:"up"`
	Type                   string   `json:"type"`
	BridgeMembers          []string `json:"bridge-members"`
	RecPackets             uint64   `json:"recPackets"`
	RecErrors              uint64   `json:"recErrors"`
	RecDrops               uint64   `json:"recDrops"`
	SentBytes              uint64   `json:"sentBytes"`
	SentPackets            uint64   `json:"sentPackets"`
	LinkSpeed              uint64   `json:"link_speed"`
	SentErrors             uint64   `json:"sentErrors"`
	SentDrops              uint64   `json:"sentDrops"`
	CarrierChanges         uint64   `json:"carrierChanges"`
	Carrier                uint64   `json:"carrier"`
	CarrierUpCount         uint64   `json:"carrier_up_count"`
	CarrierDownCount       uint64   `json:"carrier_down_count"`
	Present                bool     `json:"present"`
	External               bool     `json:"external"`
	Ipv6                   bool     `json:"ipv6"`
	Macs                   uint64   `json:"macs"`
	Mtu                    uint64   `json:"mtu"`
	LinkSupported          []string `json:"link-supported"`
	LinkAdvertising        []string `json:"link-advertising"`
	LinkPartnerAdvertising []string `json:"link-partner-advertising"`
	Multicast              bool     `json:"multicast"`
	NegotiationType        string   `json:"negotiation_type"`
	FoundDescriptor        string   `json:"foundDescriptor"`
}

type WirelessInterface struct {
	ID              *string `json:".id,omitempty"`
	Disabled        *bool   `json:"disabled,omitempty"`
	HideSSID        *bool   `json:"hide-ssid,omitempty"`
	InterfaceType   *string `json:"interface-type,omitempty"`
	Key             *string `json:"key,omitempty"`
	MACAddress      *string `json:"mac-address,omitempty"`
	Encryption      *string `json:"encryption,omitempty"`
	MasterInterface *string `json:"master-interface,omitempty"`
	Name            *string `json:"name,omitempty"`
	Running         *bool   `json:"running,omitempty"`
	SecurityProfile *string `json:"security-profile,omitempty"`
	SSID            *string `json:"ssid,omitempty"`
	Band            *string `json:"band,omitempty"`
}

type Collector struct {
	Interface         []Interface          `json:"interface,omitempty"`
	Ping              []Ping               `json:"ping,omitempty"`
	System            System               `json:"system,omitempty"`
	Wap               []Wap                `json:"wap,omitempty"`
	Tcp               TcpCollector         `json:"tcp,omitempty"`
	Sensor            Sensor               `json:"sensor,omitempty"`
	Gauge             []Gauge              `json:"gauge,omitempty"`
	Counter           []Counter            `json:"counter,omitempty"`
	MGauge            []MGauge             `json:"mgauge,omitempty"`
	MCounter          []MCounter           `json:"mcounter,omitempty"`
	LocationCollector LocationCollector    `json:"location,omitempty"`
	WifiAreaScanning  []WifiAreaScanning   `json:"wifi_area_scanning,omitempty"`
	Leases            map[string]Lease     `json:"leases,omitempty"`
	Spectralscan      json.RawMessage      `json:"spectralscan,omitempty"`
	Arptable          map[string]ArpDevice `json:"arptable,omitempty"`
	Topos             []Topo               `json:"topos,omitempty"`
	Gps               Gps                  `json:"gps,omitempty"`
}

type Wap struct {
	Interface       string    `json:"interface"`
	Ssid            string    `json:"ssid,omitempty"`
	Stations        []Station `json:"stations"`
	Signal0         float64   `json:"signal0"`
	Signal1         float64   `json:"signal1"`
	Signal2         float64   `json:"signal2"`
	Signal3         float64   `json:"signal3"`
	Dhcp            Dhcp      `json:"dhcp"`
	Noise           float64   `json:"noise"`
	FoundDescriptor string
	Key             string `json:"key"`
	Keytypes        string `json:"keytypes"`
	// added for skynet devices and we need to empliment it with our old openwrt agent
	Signal    float64 `json:"signal,omitempty"`
	Channel   int     `json:"channel,omitempty"`
	Frequency int     `json:"frequency,omitempty"`
	TxPower   int     `json:"txPower,omitempty"`
	Quality   int     `json:"quality,omitempty"`
	Rate      int     `json:"rate,omitempty"`
	Chutil    int     `json:"chutil,omitempty"`
	BandWidth int     `json:"bandWidth,omitempty"`
	Mac       string  `json:"mac,omitempty"`
	Mode      string  `json:"mode,omitempty"`
	Phy       string  `json:"phy,omitempty"`
}

type TcpCollector struct {
	UniqueIps         uint64 `json:"uniqueIps"`
	SlowedPairPackets uint64 `json:"slowedPairPackets"`
	Cwr               uint64 `json:"cwr"` // Congestion Window Reduced
	Ece               uint64 `json:"ece"` // ECN Echo
	Rst               uint64 `json:"rst"` // Reset
	Syn               uint64 `json:"syn"` // Synchronize
	Urg               uint64 `json:"urg"` // Urgent
}

type Alt struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Baro struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Batt struct {
	Name     string  `json:"name"`
	Charge   float64 `json:"charge"`
	Voltage  float64 `json:"voltage"`
	Amperage float64 `json:"amperage"`
	Temp     float64 `json:"temp"`
}

type Spl struct {
	Name      string   `json:"name"`
	Duration  float64  `json:"duration"`
	RangeSize uint64   `json:"rangeSize"`
	Counts    []uint64 `json:"counts"`
}

type Prox struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Camera struct {
	Name        string  `json:"name"`
	MotionCount uint64  `json:"motionCount"`
	BodyCount   uint64  `json:"bodyCount"`
	AnimalCount uint64  `json:"animalCount"`
	HandCount   uint64  `json:"handCount"`
	FaceCount   uint64  `json:"faceCount"`
	AvgDistance float64 `json:"avgDistance"`
	LightPct    float64 `json:"lightPct"`
}

type Env struct {
	Name     string  `json:"name"`
	Pressure float64 `json:"pressure"`
	Humidity float64 `json:"humidity"`
	Temp     float64 `json:"temp"`
	Airflow  float64 `json:"airflow"`
}

type Sensor struct {
	Alt    []Alt    `json:"alt"`
	Baro   []Baro   `json:"baro"`
	Batt   []Batt   `json:"batt"`
	Spl    []Spl    `json:"spl"`
	Prox   []Prox   `json:"prox"`
	Camera []Camera `json:"camera"`
	Env    []Env    `json:"env"`
}

type LocationCollector struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Gauge struct {
	Name  string  `json:"name"`
	Point float64 `json:"point"`
}

type Counter struct {
	Name  string `json:"name"`
	Point uint64 `json:"point"`
}

type MGauge struct {
	Name       string    `json:"name"`
	PointNames []string  `json:"pointNames"`
	Points     []float64 `json:"points"`
}

type MCounter struct {
	Name       string   `json:"name"`
	PointNames []string `json:"pointNames"`
	Points     []uint64 `json:"points"`
}

type Gps struct {
	Height    float32 `json:"height,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
	Latitude  float32 `json:"latitude,omitempty"`
}

type ArpDevice struct {
	Name    string `json:"name,omitempty"`
	Device  string `json:"device,omitempty"`
	Clients []struct {
		Ip     string `json:"ip,omitempty"`
		Mac    string `json:"mac,omitempty"`
		Flags  string `json:"flags,omitempty"`
		Hwtype string `json:"hwtype,omitempty"`
		Mask   string `json:"mask,omitempty"`
	} `json:"clients,omitempty"`
}

type Topo struct {
	Mac      string `json:"mac,omitempty"`
	PMac     string `json:"pMac,omitempty"`
	Hops     int    `json:"hops,omitempty"`
	Ip       string `json:"ip,omitempty"`
	Backhaul string `json:"backhaul,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Lease struct {
	Timestamp uint64 `json:"timestamp,omitempty"`
	Mac       string `json:"mac,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	Duid      string `json:"duid,omitempty"`
}

type WifiAreaScanning struct {
	// this for capuring scanned stations
	RadioScaning struct {
		Device    string `json:"device,ommitempty"`
		Frequency uint64 `json:"frequency,ommitempty"`
	} `json:"radio_scaning,ommitempty"`
	FoundStations []FoundStations `json:"found_stations,ommitempty"`
}

type FoundStations struct {
	// back to wifi area scanning
	QualityMax        uint64     `json:"quality_max"`
	Ssid              string     `json:"ssid"`
	Encryption        Encryption `json:"encryption"`
	Channel           int        `json:"channel,omitempty"`
	EstimatedDistance float64    `json:"estimated_distance,omitempty"`
	Quality           uint64     `json:"quality,omitempty"`
	Signal            int        `json:"signal,omitempty"`
	Mode              string     `json:"mode,omitempty"`
	// {"quality_max":94,"ssid":"Atour","encryption":{"enabled":true,"wpa":[1,2],"ciphers":["ccmp"],"authentication":["psk"]}
}

type Encryption struct {
	// stations encryption nothing to do with device but with the scanned stations
	Enabled        bool     `json:"enabled"`
	Ciphers        []string `json:"ciphers,omitempty"`
	Wep            []string `json:"wep,omitempty"`
	Wap            []int    `json:"wap,omitempty"`
	Authentication []string `json:"authentication,omitempty"`
}

type System struct {
	Load   Load   `json:"load"`
	Memory Memory `json:"memory"`
	Disks  []Disk `json:"disks"`
}

type Load struct {
	One          float64 `json:"one"`
	Five         float64 `json:"five"`
	Fifteen      float64 `json:"fifteen"`
	ProcessCount uint64  `json:"processCount"`
}

type Disk struct {
	Mount string `json:"mount"`
	Used  uint64 `json:"used"`
	Avail uint64 `json:"avail"`
}

type Memory struct {
	Total   uint64 `json:"total"`
	Free    uint64 `json:"free"`
	Buffers uint64 `json:"buffers"`
	Cache   uint64 `json:"cache"`
}

type Station struct {
	Mac             string  `json:"mac"`
	Info            string  `json:"info"`
	Rssi            float64 `json:"rssi"`
	RecBytes        uint64  `json:"recBytes"`
	SentBytes       uint64  `json:"sentBytes"`
	Ccq             float64 `json:"ccq"`
	Noise           float64 `json:"noise"`
	Signal0         float64 `json:"signal0"`
	Signal1         float64 `json:"signal1"`
	Signal2         float64 `json:"signal2"`
	Dhcp            Lease   `json:"dhcp,omitempty"`
	Signal3         float64 `json:"signal3"`
	ExpectedRate    uint64  `json:"expectedRate"`
	AssocTime       uint64  `json:"assocTime"`
	BeaconLoss      uint64  `json:"beaconLoss"`
	FoundDescriptor string
}

type Dhcp struct {
	Ip        string  `json:"ip,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
	Mac       string  `json:"mac,omitempty"`
	Hostname  string  `json:"hostname,omitempty"`
	Duid      string  `json:"duid,omitempty"`
}
type Ping struct {
	Host   string  `json:"host"`
	AvgRtt float64 `json:"avgRtt"`
	MinRtt float64 `json:"minRtt"`
	MaxRtt float64 `json:"maxRtt"`
	Loss   float64 `json:"loss"`
}
