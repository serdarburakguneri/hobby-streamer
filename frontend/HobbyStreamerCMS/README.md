# Hobby Streamer CMS – React Native (Web)

This is a simple CMS frontend built with React Native (via Expo), mainly for experimenting with asset management flows. It runs in the browser using Expo for Web — handy for local testing and quick iterations.

## Features

- Paginated asset list view  
- Asset creation and edit forms  
- File upload support (via presigned S3 URLs)  
- Basic auth integration (Keycloak)  
- Clean UI built with React Navigation

---

## Prerequisites

Before running the app, make sure you have:

- **Node.js** (v18 or higher)
- **npm** or **yarn**
- **Expo CLI** (installed as part of dependencies)

---

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Start the dev server (for web):

```bash
npm run web
```

3. Then open your browser at:

```
http://localhost:8001
```

---

## Project Structure

```
src/
├── components/     # Shared/reusable UI pieces
├── screens/        # Page-level components
├── services/       # API calls and backend integration
├── types/          # TypeScript definitions
└── utils/          # Helpers and utilities
```

---

## Backend Integration

The CMS connects to locally running backend services:

- **Asset Manager** – `http://localhost:8082`  
- **Auth Service** – `http://localhost:8080`

These are expected to be running via Docker Compose (see main project root).

---

## Development Notes

### Hot Reloading

This project uses Expo with hot reload, so changes to most frontend code should reflect instantly:

```bash
npm run web       # Launch web version with hot reload
npm run ios       # Run on iOS simulator (if available)
npm run android   # Run on Android emulator
```

What **updates automatically**:
- UI styling changes
- Component structure (JSX)
- Props and local state

What **needs restart**:
- Native modules
- `expo.config.js`
- Dependency changes

---

### Manual Refresh (if needed)

If things get stuck, try:

- Browser refresh (`F5` or `Cmd+R`)
- Terminal shortcut (`r` when Metro is running)
- Mobile: shake device and select "Reload"

---

### Adding New Screens

1. Create a file under `src/screens/`
2. Register the screen in the main navigator (in `App.tsx`)
3. Update type definitions if needed (in `types/navigation.ts`)

---

### API Integration

All API logic lives in `src/services/api.ts`:

- `assetService`: handles upload, fetch, update, delete
- `authService`: login and token management

No global state yet — this is intentionally minimal and modular while exploring things.

---

> ⚠️ This project is part of the broader Hobby Streamer playground. It's not production-ready, but serves as a testing ground for UI and API patterns during local dev.