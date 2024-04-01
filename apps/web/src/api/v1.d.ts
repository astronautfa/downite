/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */


export interface paths {
  "/api/v1/torrent": {
    /**
     * GetTorrents
     * @description controller: downite/handlers.GetTorrents
     */
    get: operations["GET /api/v1/torrent:GetTorrents"];
  };
  "/api/v1/torrent-meta": {
    /**
     * GetTorrentMeta
     * @description controller: downite/handlers.GetTorrentMeta
     */
    post: operations["POST /api/v1/torrent-meta:GetTorrentMeta"];
  };
}

export type webhooks = Record<string, never>;

export interface components {
  schemas: {
    FileTree: {
      Dir?: {
        [key: string]: components["schemas"]["FileTree"];
      };
      File?: {
        /** Format: int64 */
        Length?: number;
        PiecesRoot?: string;
      };
    };
    GetTorrentMetaReq: {
      /** Format: byte */
      file?: string;
      magnet?: string;
    };
    GetTorrentMetaRes: {
      files?: {
          /** Format: int64 */
          Length?: number;
          Path?: string[];
        }[];
      name?: string;
      /** Format: int64 */
      totalSize?: number;
    };
    Torrent: {
        /** Format: int64 */
        addedOn?: number;
        /** Format: float */
        availability?: number;
        category?: string;
        downloadDir?: string;
        downloadPath?: string;
        downloadSpeed?: number;
        eta?: number;
        files?: {
          Dir?: {
            [key: string]: components["schemas"]["FileTree"];
          };
          File?: {
            /** Format: int64 */
            Length?: number;
            PiecesRoot?: string;
          };
        };
        infoHash?: string;
        name?: string;
        peers?: {
          [key: string]: {
            Addr?: unknown;
            Id?: unknown;
            Source?: string;
            SupportsEncryption?: boolean;
            Trusted?: boolean;
          };
        };
        peersCount?: number;
        pieceProgress?: {
            DownloadedByteCount?: number;
            Index?: number;
            Length?: number;
          }[];
        /** Format: float */
        progress?: number;
        /** Format: float */
        ratio?: number;
        seeds?: number;
        status?: number;
        tags?: string[];
        totalSize?: number;
        uploadSpeed?: number;
      }[];
  };
  responses: never;
  parameters: never;
  requestBodies: {
    /** @description Request body for handlers.GetTorrentMetaReq */
    GetTorrentMetaReq: {
      content: {
        "application/json": components["schemas"]["GetTorrentMetaReq"];
      };
    };
  };
  headers: never;
  pathItems: never;
}

export type $defs = Record<string, never>;

export type external = Record<string, never>;

export interface operations {

  /**
   * GetTorrents
   * @description controller: downite/handlers.GetTorrents
   */
  "GET /api/v1/torrent:GetTorrents": {
    responses: {
      /** @description OK */
      200: {
        content: {
          "application/json": components["schemas"]["Torrent"];
        };
      };
      default: {
        content: never;
      };
    };
  };
  /**
   * GetTorrentMeta
   * @description controller: downite/handlers.GetTorrentMeta
   */
  "POST /api/v1/torrent-meta:GetTorrentMeta": {
    requestBody: components["requestBodies"]["GetTorrentMetaReq"];
    responses: {
      /** @description OK */
      200: {
        content: {
          "application/json": components["schemas"]["GetTorrentMetaRes"];
        };
      };
      default: {
        content: never;
      };
    };
  };
}