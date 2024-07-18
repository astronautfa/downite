/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as TorrentsImport } from './routes/torrents'
import { Route as TorrentImport } from './routes/torrent'
import { Route as SettingsImport } from './routes/settings'
import { Route as HelpImport } from './routes/help'
import { Route as DownloadsImport } from './routes/downloads'
import { Route as AccountImport } from './routes/account'
import { Route as IndexImport } from './routes/index'
import { Route as TorrentsIndexImport } from './routes/torrents/index'
import { Route as DownloadsIndexImport } from './routes/downloads/index'
import { Route as TorrentInfohashImport } from './routes/torrent/$infohash'
import { Route as DownloadIdImport } from './routes/download/$id'

// Create/Update Routes

const TorrentsRoute = TorrentsImport.update({
  path: '/torrents',
  getParentRoute: () => rootRoute,
} as any)

const TorrentRoute = TorrentImport.update({
  path: '/torrent',
  getParentRoute: () => rootRoute,
} as any)

const SettingsRoute = SettingsImport.update({
  path: '/settings',
  getParentRoute: () => rootRoute,
} as any)

const HelpRoute = HelpImport.update({
  path: '/help',
  getParentRoute: () => rootRoute,
} as any)

const DownloadsRoute = DownloadsImport.update({
  path: '/downloads',
  getParentRoute: () => rootRoute,
} as any)

const AccountRoute = AccountImport.update({
  path: '/account',
  getParentRoute: () => rootRoute,
} as any)

const IndexRoute = IndexImport.update({
  path: '/',
  getParentRoute: () => rootRoute,
} as any)

const TorrentsIndexRoute = TorrentsIndexImport.update({
  path: '/',
  getParentRoute: () => TorrentsRoute,
} as any)

const DownloadsIndexRoute = DownloadsIndexImport.update({
  path: '/',
  getParentRoute: () => DownloadsRoute,
} as any)

const TorrentInfohashRoute = TorrentInfohashImport.update({
  path: '/$infohash',
  getParentRoute: () => TorrentRoute,
} as any)

const DownloadIdRoute = DownloadIdImport.update({
  path: '/download/$id',
  getParentRoute: () => rootRoute,
} as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexImport
      parentRoute: typeof rootRoute
    }
    '/account': {
      id: '/account'
      path: '/account'
      fullPath: '/account'
      preLoaderRoute: typeof AccountImport
      parentRoute: typeof rootRoute
    }
    '/downloads': {
      id: '/downloads'
      path: '/downloads'
      fullPath: '/downloads'
      preLoaderRoute: typeof DownloadsImport
      parentRoute: typeof rootRoute
    }
    '/help': {
      id: '/help'
      path: '/help'
      fullPath: '/help'
      preLoaderRoute: typeof HelpImport
      parentRoute: typeof rootRoute
    }
    '/settings': {
      id: '/settings'
      path: '/settings'
      fullPath: '/settings'
      preLoaderRoute: typeof SettingsImport
      parentRoute: typeof rootRoute
    }
    '/torrent': {
      id: '/torrent'
      path: '/torrent'
      fullPath: '/torrent'
      preLoaderRoute: typeof TorrentImport
      parentRoute: typeof rootRoute
    }
    '/torrents': {
      id: '/torrents'
      path: '/torrents'
      fullPath: '/torrents'
      preLoaderRoute: typeof TorrentsImport
      parentRoute: typeof rootRoute
    }
    '/download/$id': {
      id: '/download/$id'
      path: '/download/$id'
      fullPath: '/download/$id'
      preLoaderRoute: typeof DownloadIdImport
      parentRoute: typeof rootRoute
    }
    '/torrent/$infohash': {
      id: '/torrent/$infohash'
      path: '/$infohash'
      fullPath: '/torrent/$infohash'
      preLoaderRoute: typeof TorrentInfohashImport
      parentRoute: typeof TorrentImport
    }
    '/downloads/': {
      id: '/downloads/'
      path: '/'
      fullPath: '/downloads/'
      preLoaderRoute: typeof DownloadsIndexImport
      parentRoute: typeof DownloadsImport
    }
    '/torrents/': {
      id: '/torrents/'
      path: '/'
      fullPath: '/torrents/'
      preLoaderRoute: typeof TorrentsIndexImport
      parentRoute: typeof TorrentsImport
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren({
  IndexRoute,
  AccountRoute,
  DownloadsRoute: DownloadsRoute.addChildren({ DownloadsIndexRoute }),
  HelpRoute,
  SettingsRoute,
  TorrentRoute: TorrentRoute.addChildren({ TorrentInfohashRoute }),
  TorrentsRoute: TorrentsRoute.addChildren({ TorrentsIndexRoute }),
  DownloadIdRoute,
})

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/account",
        "/downloads",
        "/help",
        "/settings",
        "/torrent",
        "/torrents",
        "/download/$id"
      ]
    },
    "/": {
      "filePath": "index.tsx"
    },
    "/account": {
      "filePath": "account.tsx"
    },
    "/downloads": {
      "filePath": "downloads.tsx",
      "children": [
        "/downloads/"
      ]
    },
    "/help": {
      "filePath": "help.tsx"
    },
    "/settings": {
      "filePath": "settings.tsx"
    },
    "/torrent": {
      "filePath": "torrent.tsx",
      "children": [
        "/torrent/$infohash"
      ]
    },
    "/torrents": {
      "filePath": "torrents.tsx",
      "children": [
        "/torrents/"
      ]
    },
    "/download/$id": {
      "filePath": "download/$id.tsx"
    },
    "/torrent/$infohash": {
      "filePath": "torrent/$infohash.tsx",
      "parent": "/torrent"
    },
    "/downloads/": {
      "filePath": "downloads/index.tsx",
      "parent": "/downloads"
    },
    "/torrents/": {
      "filePath": "torrents/index.tsx",
      "parent": "/torrents"
    }
  }
}
ROUTE_MANIFEST_END */
