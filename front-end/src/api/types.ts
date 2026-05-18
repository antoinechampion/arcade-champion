export type Platform = 'Steam' | 'Fightcade' | 'MAME'

export interface Game {
  id: string
  title: string
  platform: Platform
  releaseYear: number
  developer: string
  coverFilename: string
  bannerFilename: string
  appId: string
}

export interface GameInput {
  title: string
  platform: Platform
  releaseYear: number
  developer: string
  appId: string
  cover?: Blob
  banner?: Blob
}

export interface PlatformSearchResult {
  name: string
  platformId: string
}

export interface Settings {
  fightcadeUsername: string
  fightcadePassword: string
  fightcadeCookie: string
  mamePath: string
  steamPath: string
}
