BEGIN {
    FS = "[ \t]+"
    RS = ""
    print "["
    first_station = 1
}

function parse_station_line(line) {
    split(line, parts)
    station["addr"] = parts[1]
    station["aid"] = parts[2]
    station["chan"] = parts[3]
    station["txrate"] = parts[4]
    station["rxrate"] = parts[5]
    station["rssi"] = parts[6]
    station["minrssi"] = parts[7]
    station["maxrssi"] = parts[8]
    station["idle"] = parts[9]
    station["txseq"] = parts[10]
    station["rxseq"] = parts[11]
    station["caps"] = parts[12]
    station["xcaps"] = parts[13]
    station["acaps"] = parts[14]
    station["erp"] = parts[15]
    station["state"] = parts[16]
    station["maxrate(dot11)"] = parts[17]
    station["htcaps"] = parts[18]
    station["vhtcaps"] = parts[19]
    station["assoctime"] = parts[20]
    station["ies"] = parts[21]
}

function parse_additional_info(line) {
    if (line ~ /RSSI is combined over chains in dBm/) {
        station["rssi_combined"] = "true"
    } else if (line ~ /Minimum Tx Power/) {
        station["min_tx_power"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Maximum Tx Power/) {
        station["max_tx_power"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /HT Capability/) {
        station["ht_capability"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /VHT Capability/) {
        station["vht_capability"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /MU capable/) {
        station["mu_capable"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /SNR/) {
        station["snr"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Operating band/) {
        station["operating_band"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Current Operating class/) {
        station["current_operating_class"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Supported Operating classes/) {
        station["supported_operating_classes"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Supported Rates\(Mbps\)/) {
        station["supported_rates"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Max STA phymode/) {
        station["max_sta_phymode"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /MLO/) {
        station["mlo"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /MLD Addr/) {
        station["mld_addr"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Num Partner links/) {
        station["num_partner_links"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /Partner link/) {
        if (!("partner_links" in station)) {
            station["partner_links"] = ""
        }
        station["partner_links"] = station["partner_links"] " " substr(line, index(line, ":") + 2)
    } else if (line ~ /EMLSR capable/) {
        station["emlsr_capable"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /EMLMR capable/) {
        station["emlmr_capable"] = substr(line, index(line, ":") + 2)
    } else if (line ~ /STR capable/) {
        station["str_capable"] = substr(line, index(line, ":") + 2)
    }
}

{
    split($0, lines, "\n")
    for (i in lines) {
        line = lines[i]
        if (match(line, /^[0-9a-f]{2}(:[0-9a-f]{2}){5}/)) {
            if (length(station) > 0) {
                print_station()
            }
            delete station
            parse_station_line(line)
        } else if (length(station) > 0) {
            parse_additional_info(line)
        }
    }
    if (length(station) > 0) {
        print_station()
    }
}

function print_station() {
    if (!first_station) {
        printf(",\n")
    }
    first_station = 0
    printf("{")
    first_field = 1
    for (key in station) {
        if (!first_field) {
            printf(",")
        }
        first_field = 0
        printf("\"" key "\": \"" station[key] "\"")
    }
    printf("}")
}

END {
    print "\n]"
}