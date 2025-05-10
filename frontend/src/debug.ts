import { devtools } from '@vue/devtools'
import { GetDebugConfig } from '../wailsjs/go/controller/Controller'

const debugConfig = await GetDebugConfig()

if (debugConfig.VueDevTools.Host !== '') {
    let host = debugConfig.VueDevTools.Host
    if(!host.startsWith('http')) {
        host = 'http://' + host
    }
    console.log('Connecting to Vue-Devtools', host, debugConfig.VueDevTools.Port)

    devtools.connect(host, debugConfig.VueDevTools.Port).then(() => {
        console.log('Vue-Devtools connected!')
    })
} else {
    console.debug('Vue-Devtools not enabled')
}
