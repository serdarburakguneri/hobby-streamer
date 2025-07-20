module.exports = function(api) {
  api.cache(true);
  return {
    presets: [
      ['babel-preset-expo', {
        web: {
          unstable_transformProfile: 'hermes-stable'
        }
      }]
    ],
    plugins: [
      // Removed react-native-reanimated/plugin as it's not installed
    ],
  };
}; 