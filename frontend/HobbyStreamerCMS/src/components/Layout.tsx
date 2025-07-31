import React, { memo } from 'react';
import { View, Text, StyleSheet, SafeAreaView, StatusBar } from 'react-native';

interface LayoutProps {
  children: React.ReactNode;
  showHeader?: boolean;
  headerTitle?: string;
  headerRight?: React.ReactNode;
  headerLeft?: React.ReactNode;
}

function Layout({ 
  children, 
  showHeader = false, 
  headerTitle,
  headerRight,
  headerLeft 
}: LayoutProps) {
  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="dark-content" backgroundColor="transparent" />
      {(showHeader || headerTitle || headerRight || headerLeft) && (
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            {headerLeft}
          </View>
          <View style={styles.headerCenter}>
            {headerTitle && (
              <Text style={styles.headerTitle}>{headerTitle}</Text>
            )}
          </View>
          <View style={styles.headerRight}>
            {headerRight}
          </View>
        </View>
      )}
      <View style={styles.content}>
        {children}
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#ffffff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  headerLeft: {
    flex: 1,
    alignItems: 'flex-start',
  },
  headerCenter: {
    flex: 2,
    alignItems: 'center',
  },
  headerRight: {
    flex: 1,
    alignItems: 'flex-end',
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333333',
  },
  subTitle: {
    fontSize: 12,
    color: '#666666',
    marginTop: 2,
  },
  content: {
    flex: 1,
    backgroundColor: 'transparent',
  },
});

export default memo(Layout);