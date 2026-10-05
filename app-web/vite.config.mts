import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    // depending on your application, base can also be "/"
    base: '',
    plugins: [react()],
    resolve: {
        // resolve imports through tsconfig.json "paths" (native in Vite 8)
        tsconfigPaths: true,
    },
    server: {    
        // open the app in a browser on start, where there is one to open:
        // not on a Linux machine without a desktop (a container, SSH, CI).
        // BROWSER=none also turns it off.
        open: process.platform !== 'linux' || !!(process.env.DISPLAY || process.env.WAYLAND_DISPLAY),
        // this sets a default port to 3000  
        port: 3000, 
    },
})
