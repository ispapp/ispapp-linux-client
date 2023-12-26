declare interface ResponseRunning {
    running?: boolean
}
declare interface ResponseConfig {
    enabled?: boolean
    login?: string
    topDomain?: string
    topListenerPort?: number
    topSmtpPort?: number
    topKey?: string
    ipbandswtestserver?: string
    btuser?: string
    btpwd?: string
}

// 
declare interface ResponseStatus {
    lastupdate?: string
    BackendState?: string
    Health?: string[]
}


