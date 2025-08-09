function fn() {
    var env = karate.env; 
    karate.log('karate.env system property was:', env);
    if (!env) {
        env = 'local';
    }
    var apiGatewayId = 'unknown';
    if (env == 'local') {
        var fs = Java.type('java.nio.file.Files');
        var path = Java.type('java.nio.file.Paths');
        var apiGatewayIdPath = path.get('../local/.api-gateway-id');
        if (fs.exists(apiGatewayIdPath)) {
            apiGatewayId = fs.readString(apiGatewayIdPath).trim();
            karate.log('API Gateway ID loaded from ../local/api-gateway-id:', apiGatewayId);
        } else {
            karate.log('API Gateway ID file ../local/api-gateway-id not found.');
        }
    }
    
    
    var config = {
        env: env,
        assetManagerUrl: 'http://localhost:8082/graphql',
        authServiceUrl: 'http://localhost:8080',
        streamingApiUrl: 'http://localhost:8084',
        nginxUrl: 'http://localhost:8083',
        neo4jUrl: 'http://localhost:7474',
        localstackUrl: 'http://localhost:4566',
        keycloakUrl: 'http://localhost:9090',
        apiGatewayUrl: 'http://localhost:4566/_aws/execute-api/' + apiGatewayId + '/dev'
    };
    
    if (env == 'local') {       
        config.assetManagerUrl = 'http://localhost:8082/graphql';
        config.authServiceUrl = 'http://localhost:8080';
        config.streamingApiUrl = 'http://localhost:8084';
        config.nginxUrl = 'http://localhost:8083';
        config.keycloakUrl = 'http://localhost:9090';
        config.apiGatewayUrl = 'http://localhost:4566/_aws/execute-api/' + apiGatewayId + '/dev';
    } else if (env == 'dev') {
       
    } else if (env == 'e2e') {
        
    }
    
    karate.configure('connectTimeout', 30000);
    karate.configure('readTimeout', 30000);
    karate.configure('retry', { count: 3, interval: 1000 });
    karate.configure('logPrettyRequest', true);
    karate.configure('logPrettyResponse', true);
    karate.configure('printEnabled', true);
    
    return config;
} 