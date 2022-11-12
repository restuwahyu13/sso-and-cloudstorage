module.exports = {
	apps: [
		{
			name: 'golang-app',
			script: 'main.js',
			watch: false,
			// env: {
			// 	PORT: process.env.PORT,
			// 	NODE_ENV: process.env.GO_ENV
			// },
			exec_mode: 'cluster',
			instances: 'max',
			max_memory_restart: '512M',
			listen_timeout: 3000,
			kill_timeout: 6000,
			combine_logs: true
		}
	]
}
