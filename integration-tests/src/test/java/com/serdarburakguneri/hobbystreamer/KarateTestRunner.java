package com.serdarburakguneri.hobbystreamer;

import com.intuit.karate.junit5.Karate;

class KarateTestRunner {

    @Karate.Test
    Karate testAssetManager() {       
        return Karate.run("classpath:features/asset-manager/suites/");        
    }

    @Karate.Test
    Karate testAssetManagerBucketSuite() {       
        return Karate.run("classpath:features/asset-manager/suites/bucket-suite.feature");        
    }

    @Karate.Test
    Karate testAssetManagerAssetSuite() {       
        return Karate.run("classpath:features/asset-manager/suites/asset-suite.feature");        
    }
    
    @Karate.Test
    Karate testDataPopulate() {       
        return Karate.run("classpath:features/asset-manager/suites/data-populate.feature");        
    }
    
    @Karate.Test
    Karate testAuthService() {       
        return Karate.run("classpath:features/auth-service/auth-service.feature");        
    }
    
    @Karate.Test
    Karate testStreamingApi() {       
        return Karate.run("classpath:features/streaming-api/suites/");        
    }
    
}