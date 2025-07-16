# Hobby Streamer CMS - React Native Web App

React Native CMS interface for managing assets, built with Expo and running on the web.

## Features

- Asset listing with pagination
- Asset creation and editing
- File upload functionality
- Authentication integration
- Modern UI with React Navigation

## Prerequisites

- Node.js >= 18
- npm or yarn
- Expo CLI (installed automatically with dependencies)

## Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server (web):
```bash
npm run web
```

3. Open your browser and go to:
```
http://localhost:19006
```

## Project Structure

```
src/
├── components/     # Reusable UI components
├── screens/        # Screen components
├── services/       # API services
├── types/          # TypeScript type definitions
└── utils/          # Utility functions
```

## Backend Integration

The app connects to the following backend services:

- Asset Manager: `http://localhost:8082` - Asset CRUD operations
- Auth Service: `http://localhost:8080` - Authentication

## Development

### Hot Reload Development
This project uses Expo with hot reloading enabled. Most UI changes will update instantly without restarting:

```bash
# Start development server with hot reload
npm run web

# For iOS simulator
npm run ios

# For Android emulator  
npm run android
```

**What updates automatically:**
- Style changes
- Component JSX structure  
- Props changes
- Most state management changes

**What requires restart:**
- Native module changes
- Expo config changes
- Package.json dependencies

### Manual Reload Options
If hot reload isn't working:
- Browser refresh (F5 or Cmd+R)
- Press `r` in terminal
- Shake device and tap "Reload" (mobile)

### Adding New Screens
1. Create a new screen component in `src/screens/`
2. Add navigation routes in the main App component
3. Update the navigation types if needed

### API Integration
- API services are defined in `src/services/api.ts`
- Use the `assetService` for asset operations
- Use the `authService` for authentication



