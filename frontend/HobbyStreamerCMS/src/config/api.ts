const getEnvVar = (key: string, defaultValue: string): string => {
  if (typeof window !== 'undefined' && (window as any).__ENV__) {
    return (window as any).__ENV__[key] || defaultValue;
  }
  return process.env[key] || defaultValue;
};

export const API_CONFIG = {
  API_GATEWAY_ID: getEnvVar('REACT_APP_API_GATEWAY_ID', '88cyvd7wrl'),
  API_GATEWAY_BASE_URL: getEnvVar('REACT_APP_API_GATEWAY_BASE_URL', 'http://localhost:4566/_aws/execute-api/88cyvd7wrl/dev'),
  AUTH_BASE_URL: getEnvVar('REACT_APP_AUTH_BASE_URL', 'http://localhost:4566/_aws/execute-api/88cyvd7wrl/dev/auth'),
  GRAPHQL_BASE_URL: getEnvVar('REACT_APP_GRAPHQL_BASE_URL', 'http://localhost:4566/_aws/execute-api/88cyvd7wrl/dev/graphql'),
  LOCALSTACK_BASE_URL: getEnvVar('REACT_APP_LOCALSTACK_BASE_URL', 'http://localhost:4566'),
  CDN_BASE_URL: getEnvVar('REACT_APP_CDN_BASE_URL', 'http://localhost:8083/cdn')
}; 