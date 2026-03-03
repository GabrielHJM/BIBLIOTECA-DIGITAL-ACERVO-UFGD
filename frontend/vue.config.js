const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
	transpileDependencies: true,
	devServer: {
		port: 8081,
		host: '0.0.0.0',
		allowedHosts: 'all',
		open: false,
		client: {
			webSocketURL: 'ws://0.0.0.0:8081/ws',
		}
	}
})
