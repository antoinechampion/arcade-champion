export type Platform = 'Steam' | 'Fightcade' | 'MAME'

export interface SteamConfig {
  appId: string
}

export interface FightcadeConfig {
  romName: string
}

export interface MameConfig {
  romPath: string
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
