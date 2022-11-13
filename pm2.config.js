module.exports = {
	apps: [
		{
			name: 'go-app',
			namespace: 'golang',
			script: 'main.js',
			watch: false,
			autorestart: true,
			exec_mode: 'cluster',
			instances: 'max',
			max_memory_restart: '512M',
			listen_timeout: 4000,
			kill_timeout: 6000,
			combine_logs: true
		}
	]
}
