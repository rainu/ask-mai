import { devtools } from '@vue/devtools'
import { GetApplicationConfig } from '../wailsjs/go/controller/Controller'

const appConfig = await GetApplicationConfig()

if (appConfig.Debug.VueDevTools.Host !== '') {
    let host = appConfig.Debug.VueDevTools.Host
    if(!host.startsWith('http')) {
        host = 'http://' + host
    }
    console.log('Connecting to Vue-Devtools', host, appConfig.Debug.VueDevTools.Port)

    devtools.connect(host, appConfig.Debug.VueDevTools.Port).then(() => {
        console.log('Vue-Devtools connected!')
    })
} else {
    console.debug('Vue-Devtools not enabled')
}
