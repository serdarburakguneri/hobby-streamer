import React, { useState, useEffect, useMemo } from 'react';
import { SafeAreaView, StatusBar, Text, View, StyleSheet, ActivityIndicator, TouchableOpacity, Image } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { ApolloProvider } from '@apollo/client';
import { ApolloClient, InMemoryCache, createHttpLink, from } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import LoginScreen from './src/screens/LoginScreen';
import AssetListScreen from './src/screens/AssetListScreen';
import CreateAssetScreen from './src/screens/CreateAssetScreen';
import BucketListScreen from './src/screens/BucketListScreen';
import CreateBucketScreen from './src/screens/CreateBucketScreen';

import { getAuthToken, setAuthTokens, clearAuthTokens, setLogoutCallback, validateTokenLocally, isTokenExpiringSoon, refreshTokenIfNeeded } from './src/services/api';
import { API_CONFIG } from './src/config/api';

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentScreen, setCurrentScreen] = useState<'assets' | 'createAsset' | 'buckets' | 'createBucket'>('assets');
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);


  const apolloClient = useMemo(() => {
    const httpLink = createHttpLink({
      uri: API_CONFIG.GRAPHQL_BASE_URL,
    });

    const authLink = setContext(async (_, { headers }) => {
      const token = await refreshTokenIfNeeded();
      return {
        headers: {
          ...headers,
          authorization: token ? `Bearer ${token}` : "",
        }
      };
    });

    const errorLink = onError(({ graphQLErrors, networkError, operation, forward }) => {
      if (graphQLErrors) {
        graphQLErrors.forEach(async ({ message, extensions }) => {
          console.error('GraphQL error:', message, extensions);
          if (extensions?.code === 'UNAUTHENTICATED') {
            console.log('Token expired - attempting automatic refresh');
            
            try {
              const newToken = await refreshTokenIfNeeded();
              if (newToken) {
                console.log('Token refreshed successfully, retrying request');
                return forward(operation);
                             } else {
                 console.log('Token refresh failed - logging out');
                 await clearAuthTokens();
                 setIsAuthenticated(false);
               }
             } catch (error) {
               console.error('Token refresh error:', error);
               await clearAuthTokens();
               setIsAuthenticated(false);
             }
          }
        });
      }
      if (networkError) {
        console.error('Network error:', networkError);
        console.error('Network error details:', networkError.message);
      }
    });

    return new ApolloClient({
      link: from([errorLink, authLink, httpLink]),
      cache: new InMemoryCache(),
    });
  }, [isAuthenticated]);

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        setLogoutCallback(async () => {
          console.log('Token expired - automatically logging out');
          await clearAuthTokens();
          setIsAuthenticated(false);
          setError(null);
        });

        const token = await getAuthToken();
        console.log('App startup - Token found in storage:', !!token);
        
        if (token) {
          try {
            console.log('Validating token locally...');
            const validation = validateTokenLocally(token);
            console.log('Validation response:', validation);
            
            if (validation.valid) {
              console.log('Token is valid, setting authenticated to true');
              setIsAuthenticated(true);
              setError(null);
            } else {
              console.log('Token is invalid, clearing and setting authenticated to false');
              await clearAuthTokens();
              setIsAuthenticated(false);
            }
          } catch (validationError) {
            console.log('Token validation failed with error:', validationError);
            await clearAuthTokens();
            setIsAuthenticated(false);
          }
        } else {
          console.log('No token found, setting authenticated to false');
          setIsAuthenticated(false);
        }
      } catch (err) {
        console.error('Auth check error:', err);
        setError('Failed to check authentication status');
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuthentication();
  }, []);



  const handleLoginSuccess = async (tokens: { accessToken: string; refreshToken: string }) => {
    try {
      await setAuthTokens(tokens.accessToken, tokens.refreshToken);
      setIsAuthenticated(true);
      setError(null);
      setCurrentScreen('assets');
    } catch (err) {
      setError('Failed to save authentication tokens');
      console.error('Login error:', err);
    }
  };

  const handleLogout = async () => {
    console.log('handleLogout called');
    try {
      console.log('Clearing auth tokens...');
      await clearAuthTokens();
      console.log('Auth tokens cleared, setting authenticated to false');
      setIsAuthenticated(false);
      setError(null);
      console.log('Logout completed successfully');
    } catch (err) {
      console.error('Logout error:', err);
      setError('Failed to logout');
    }
  };

  if (isLoading) {
    return (
      <SafeAreaView style={styles.container}>
        <StatusBar barStyle="dark-content" backgroundColor="#f5f5f5" />
        <View style={styles.centerContent}>
          <ActivityIndicator size="large" color="#007AFF" style={styles.loadingSpinner} />
          <Text style={styles.loadingText}>Loading Hobby Streamer CMS...</Text>
          <Text style={styles.loadingSubtext}>Checking authentication status</Text>
        </View>
      </SafeAreaView>
    );
  }

  if (error && isAuthenticated) {
    return (
      <SafeAreaView style={styles.container}>
        <StatusBar barStyle="dark-content" backgroundColor="#f5f5f5" />
        <View style={styles.centerContent}>
          <Text style={styles.errorText}>Error: {error}</Text>
          <Text style={styles.retryText}>Please refresh the page to try again.</Text>
        </View>
      </SafeAreaView>
    );
  }

  const handleCreateAsset = () => {
    setCurrentScreen('createAsset');
  };

  const handleCreateBucket = () => {
    setCurrentScreen('createBucket');
  };

  const handleBackToAssets = () => {
    setCurrentScreen('assets');
  };

  const handleBackToBuckets = () => {
    setCurrentScreen('buckets');
  };

  const handleAssetCreated = () => {
    setRefreshTrigger(prev => prev + 1);
    setCurrentScreen('assets');
  };

  const handleBucketCreated = () => {
    setRefreshTrigger(prev => prev + 1);
    setCurrentScreen('buckets');
  };

  const handleLogoutWithConfirmation = () => {
    console.log('Logout confirmation dialog shown');
    setShowLogoutConfirm(true);
  };

  const handleLogoutCancel = () => {
    console.log('Logout cancelled');
    setShowLogoutConfirm(false);
  };

  const handleLogoutConfirm = () => {
    console.log('Logout confirmed');
    setShowLogoutConfirm(false);
    handleLogout();
  };

  return (
    <ApolloProvider client={apolloClient}>
      <View style={styles.container}>
        <StatusBar barStyle="dark-content" backgroundColor="#f5f5f5" />

        

        {isAuthenticated ? (
          <>
            {currentScreen === 'assets' || currentScreen === 'buckets' ? (
              <View style={styles.mainLayout}>
                <View style={styles.sidebar}>
                  <View style={styles.sidebarHeader}>
                    <View style={styles.sidebarTitleContainer}>
                      <View style={styles.logo}>
                        <Image 
                          source={require('./assets/logo.png')} 
                          style={styles.logoImage}
                          resizeMode="contain"
                        />
                      </View>
                    </View>
                  </View>
                  
                  <TouchableOpacity
                    style={[styles.sidebarItem, currentScreen === 'assets' && styles.activeSidebarItem]}
                    onPress={() => setCurrentScreen('assets')}
                  >
                    <Ionicons 
                      name="videocam" 
                      size={24} 
                      color={currentScreen === 'assets' ? '#007AFF' : '#bdc3c7'}
                    />
                    <Text style={[styles.sidebarItemText, currentScreen === 'assets' && styles.activeSidebarItemText]}>
                      Assets
                    </Text>
                  </TouchableOpacity>
                  
                  <TouchableOpacity
                    style={[styles.sidebarItem, currentScreen === 'buckets' && styles.activeSidebarItem]}
                    onPress={() => setCurrentScreen('buckets')}
                  >
                    <Ionicons 
                      name="library" 
                      size={24} 
                      color={currentScreen === 'buckets' ? '#007AFF' : '#bdc3c7'}
                    />
                    <Text style={[styles.sidebarItemText, currentScreen === 'buckets' && styles.activeSidebarItemText]}>
                      Buckets
                    </Text>
                  </TouchableOpacity>
                  
                  {/* Future navigation items can be added here */}
                  <View style={styles.sidebarSpacer} />
                  
                  {!showLogoutConfirm ? (
                    <TouchableOpacity
                      style={styles.sidebarItem}
                      onPress={() => {
                        console.log('Logout button pressed!');
                        handleLogoutWithConfirmation();
                      }}
                      activeOpacity={0.7}
                    >
                      <Ionicons name="log-out" size={20} color="#ff3b30" />
                      <Text style={[styles.sidebarItemText, { color: '#ff3b30' }]}>
                        Logout
                      </Text>
                    </TouchableOpacity>
                  ) : (
                    <View style={styles.logoutConfirmContainer}>
                      <Text style={styles.logoutConfirmText}>
                        Logout?
                      </Text>
                      <View style={styles.logoutConfirmButtons}>
                        <TouchableOpacity
                          style={styles.logoutConfirmButton}
                          onPress={handleLogoutCancel}
                        >
                          <Text style={styles.logoutConfirmButtonText}>Cancel</Text>
                        </TouchableOpacity>
                        <TouchableOpacity
                          style={[styles.logoutConfirmButton, styles.logoutConfirmButtonDestructive]}
                          onPress={handleLogoutConfirm}
                        >
                          <Text style={[styles.logoutConfirmButtonText, styles.logoutConfirmButtonTextDestructive]}>
                            Yes
                          </Text>
                        </TouchableOpacity>
                      </View>
                    </View>
                  )}
                </View>
                
                <View style={styles.contentArea}>
                  {currentScreen === 'assets' ? (
                    <AssetListScreen 
                      onCreateAsset={handleCreateAsset} 
                      refreshTrigger={refreshTrigger}
                    />
                  ) : (
                    <BucketListScreen 
                      onCreateBucket={handleCreateBucket} 
                      refreshTrigger={refreshTrigger}
                    />
                  )}
                </View>
              </View>
            ) : currentScreen === 'createAsset' ? (
              <CreateAssetScreen onBack={handleBackToAssets} onAssetCreated={handleAssetCreated} />
            ) : currentScreen === 'createBucket' ? (
              <CreateBucketScreen onBack={handleBackToBuckets} onBucketCreated={handleBucketCreated} />
            ) : null}
          </>
        ) : (
          <LoginScreen onLoginSuccess={handleLoginSuccess} />
        )}
      </View>
    </ApolloProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#1e2328',
  },
  centerContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  loadingText: {
    fontSize: 18,
    color: '#333',
    textAlign: 'center',
  },
  errorText: {
    fontSize: 16,
    color: '#d32f2f',
    textAlign: 'center',
    marginBottom: 10,
  },
  retryText: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },

  loadingSpinner: {
    marginBottom: 20,
  },
  loadingSubtext: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  mainLayout: {
    flex: 1,
    flexDirection: 'row',
  },
  sidebar: {
    width: 200,
    backgroundColor: '#0f1419',
    borderRightWidth: 1,
    borderRightColor: '#1a1d29',
    paddingTop: 20,
  },
  sidebarHeader: {
    paddingHorizontal: 20,
    paddingVertical: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#1a1d29',
    marginBottom: 20,
  },
  sidebarTitleContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
  },
  sidebarTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#ffffff',
  },
  titleStack: {
    flexDirection: 'column',
    alignItems: 'center',
  },
  sidebarMainTitle: {
    fontSize: 14,
    fontWeight: '700',
    color: '#ffffff',
    letterSpacing: 0.8,
    fontFamily: 'System',
  },
  logo: {
    marginRight: 8,
  },
  logoImage: {
    width: 128,
    height: 128,
  },
  sidebarItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingVertical: 12,
    marginBottom: 4,
  },
  activeSidebarItem: {
    backgroundColor: '#1a1d29',
    borderLeftWidth: 4,
    borderLeftColor: '#e50914',
  },
  sidebarItemText: {
    fontSize: 16,
    fontWeight: '500',
    color: '#8b8b8b',
    marginLeft: 12,
  },
  activeSidebarItemText: {
    color: '#ffffff',
    fontWeight: '600',
  },
  sidebarSpacer: {
    flex: 1,
  },
  contentArea: {
    flex: 1,
    backgroundColor: 'transparent',
  },
  logoutConfirmContainer: {
    paddingHorizontal: 20,
    paddingVertical: 12,
    marginBottom: 4,
  },
  logoutConfirmText: {
    fontSize: 14,
    color: '#ecf0f1',
    textAlign: 'center',
    marginBottom: 8,
  },
  logoutConfirmButtons: {
    flexDirection: 'row',
    gap: 8,
  },
  logoutConfirmButton: {
    flex: 1,
    backgroundColor: '#34495e',
    paddingVertical: 6,
    paddingHorizontal: 12,
    borderRadius: 4,
    alignItems: 'center',
  },
  logoutConfirmButtonDestructive: {
    backgroundColor: '#e74c3c',
  },
  logoutConfirmButtonText: {
    fontSize: 12,
    color: '#ecf0f1',
    fontWeight: '500',
  },
  logoutConfirmButtonTextDestructive: {
    color: '#fff',
  },
});
