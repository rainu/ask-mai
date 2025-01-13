import { Log } from '../wailsjs/go/controller/Controller'

const defaultLog = console.log
const defaultTrace = console.trace
const defaultError = console.error
const defaultWarn = console.warn
const defaultInfo = console.info
const defaultDebug = console.debug

window.console.log = (...data: any[]) => {
    Log('log', JSON.stringify(data))
    defaultLog(...data)
}

window.console.trace = (...data: any[]) => {
    Log('trace', JSON.stringify(data))
    defaultTrace(...data)
}

window.console.error = (...data: any[]) => {
    Log('error', JSON.stringify(data))
    defaultError(...data)
}

window.console.warn = (...data: any[]) => {
    Log('warn', JSON.stringify(data))
    defaultWarn(...data)
}

window.console.info = (...data: any[]) => {
    Log('info', JSON.stringify(data))
    defaultInfo(...data)
}

window.console.debug = (...data: any[]) => {
    Log('debug', JSON.stringify(data))
    defaultDebug(...data)
}

