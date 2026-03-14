/**
 * Wails Web Adapter
 * Intercepts Wails calls and redirects them to the HTTP API.
 */

(function() {
    console.log('Wails Web Adapter initializing...');
    window.isWeb = true;

    // 1. Mock window.go for backend calls
    window.go = new Proxy({}, {
        get: (target, pkg) => new Proxy({}, {
            get: (target, cls) => new Proxy({}, {
                get: (target, method) => (...args) => {
                    // console.log(`Calling backend: ${pkg}.${cls}.${method}`, args);
                    return fetch(`/api/${method}`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(args)
                    }).then(async r => {
                        if (!r.ok) {
                            const err = await r.json();
                            throw new Error(err.error || 'Unknown error');
                        }
                        return r.json();
                    });
                }
            })
        })
    });

    // 2. Mock window.runtime for events and other utilities
    const eventListeners = {};
    window.runtime = {
        EventsOn: (name, callback) => {
            if (!eventListeners[name]) eventListeners[name] = [];
            eventListeners[name].push(callback);
            console.log(`Registered listener for event: ${name}`);
        },
        EventsOff: (name) => {
            delete eventListeners[name];
        },
        EventsEmit: (name, ...args) => {
            // Not implemented for web -> backend yet, but could be done via POST
            console.warn(`EventsEmit called for ${name}, but web -> backend events are not implemented yet.`);
        },
        EventsOnMultiple: (name, callback, max) => {
            if (!eventListeners[name]) eventListeners[name] = [];
            eventListeners[name].push(callback);
        },
        BrowserOpenURL: (url) => window.open(url, '_blank'),
        LogDebug: (msg) => console.debug('[Backend]', msg),
        LogInfo: (msg) => console.info('[Backend]', msg),
        LogWarning: (msg) => console.warn('[Backend]', msg),
        LogError: (msg) => console.error('[Backend]', msg),
        LogPrint: (msg) => console.log('[Backend]', msg),
        LogTrace: (msg) => console.trace('[Backend]', msg),
        LogFatal: (msg) => console.error('[Backend FATAL]', msg),
        WindowSetTitle: (title) => document.title = title,
        WindowFullscreen: () => console.log('WindowFullscreen called'),
        WindowUnfullscreen: () => console.log('WindowUnfullscreen called'),
        WindowHide: () => console.log('WindowHide called'),
        WindowShow: () => console.log('WindowShow called'),
        WindowReload: () => window.location.reload(),
        WindowReloadApp: () => window.location.reload(),
        WindowSetAlwaysOnTop: () => {},
        WindowCenter: () => {},
        WindowGetSize: () => Promise.resolve({width: window.innerWidth, height: window.innerHeight}),
        WindowSetSize: () => {},
        WindowGetPosition: () => Promise.resolve({x: 0, y: 0}),
        WindowSetPosition: () => {},
        Quit: () => console.log('Quit called'),
        Environment: () => Promise.resolve({buildType: 'web', platform: 'web'}),
    };

    // 3. Setup SSE for receiving backend events
    function connectSSE() {
        console.log('Connecting to event stream...');
        const evSource = new EventSource('/events');
        
        evSource.onmessage = (event) => {
            try {
                const { name, data } = JSON.parse(event.data);
                if (eventListeners[name]) {
                    eventListeners[name].forEach(cb => cb(data));
                }
            } catch (e) {
                console.error('Error parsing event:', e);
            }
        };

        evSource.onerror = (err) => {
            console.error('SSE connection lost, retrying in 5s...', err);
            evSource.close();
            setTimeout(connectSSE, 5000);
        };
    }

    connectSSE();
})();
