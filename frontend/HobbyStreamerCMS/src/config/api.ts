const getEnvVar = (key: string, defaultValue: string): string => {
  if (typeof window !== 'undefined' && (window as any).__ENV__) {
    return (window as any).__ENV__[key] || defaultValue;
  }
  return process.env[key] || defaultValue;
};

export const API_CONFIG = {
  API_GATEWAY_ID: getEnvVar('REACT_APP_API_GATEWAY_ID', 'vq85hkrtwg'),
  API_GATEWAY_BASE_URL: getEnvVar('REACT_APP_API_GATEWAY_BASE_URL', 'http://localhost:4566/_aws/execute-api/vq85hkrtwg/dev'),
  AUTH_BASE_URL: getEnvVar('REACT_APP_AUTH_BASE_URL', 'http://localhost:4566/_aws/execute-api/vq85hkrtwg/dev/auth'),
  GRAPHQL_BASE_URL: getEnvVar('REACT_APP_GRAPHQL_BASE_URL', 'http://localhost:4566/_aws/execute-api/vq85hkrtwg/dev/graphql'),
  TRANSCODE_BASE_URL: getEnvVar('REACT_APP_TRANSCODE_BASE_URL', 'http://localhost:4566/_aws/execute-api/vq85hkrtwg/dev/transcode'),
  LOCALSTACK_BASE_URL: getEnvVar('REACT_APP_LOCALSTACK_BASE_URL', 'http://localhost:4566')
}; 