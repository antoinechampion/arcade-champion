export type Platform = 'Steam' | 'Fightcade' | 'MAME'

export interface SteamConfig {
  appId: string
}

export interface FightcadeConfig {
  gameId: string
}

export interface MameConfig {
  driverName: string
}

export type LaunchConfig = SteamConfig | FightcadeConfig | MameConfig

export interface Game {
  id: string
  title: string
  platform: Platform
  releaseYear: number
  developer: string
  imageUrl: string
  bannerUrl?: string
  launchConfig: LaunchConfig
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
}
