module.exports = {
	entry: './src/index.js',
	output: {
		path: './html',
		filename: 'app.js'
	},
	module: {
		loaders: [{
			test: /\.js$/,
			exclude: /node_modules/,
			loader: 'babel-loader'
		}]
	},
	node: {
		fs: 'empty'
	}
};
