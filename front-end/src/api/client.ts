import type { Game, PlatformSearchResult, Platform } from './types'

const games: Game[] = [
  {
    id: '1',
    title: 'Street Fighter III: 3rd Strike',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'Capcom',
    imageUrl: 'https://i.pinimg.com/736x/a8/db/46/a8db46f121860572350471b0e4405c32.jpg',
    bannerUrl: 'https://d1lss44hh2trtw.cloudfront.net/resize?type=webp&url=https%3A%2F%2Fshacknews-www.s3.amazonaws.com%2Fassets%2Farticle%2F2024%2F07%2F29%2F3rd-strike-header_feature.jpg&width=1032&sign=XEMEx9JTW7qcjq1YKxMlyIOL40hjUM0eBKwjqW6KlMU',
    launchConfig: { gameId: 'sfiii3n' },
  },
  {
    id: '2',
    title: 'Marvel VS Capcom 2',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'Capcom',
    imageUrl: 'https://wiki.supercombo.gg/images/thumb/2/29/MVSC2_Cover_Art.jpg/300px-MVSC2_Cover_Art.jpg',
    launchConfig: { driverName: 'mvsc2' },
  },
  {
    id: '3',
    title: 'King of Fighters 98',
    platform: 'Fightcade',
    releaseYear: 1998,
    developer: 'SNK',
    imageUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1y8h.jpg',
    launchConfig: { gameId: 'kof98' },
  },
  {
    id: '4',
    title: 'Street Fighter 6',
    platform: 'Steam',
    releaseYear: 2023,
    developer: 'Capcom',
    imageUrl: 'https://upload.wikimedia.org/wikipedia/en/thumb/9/94/Street_Fighter_6_box_art.jpg/250px-Street_Fighter_6_box_art.jpg',
    launchConfig: { appId: '1364780' },
  },
  {
    id: '5',
    title: 'Street Fighter 2: Champion Edition',
    platform: 'Fightcade',
    releaseYear: 1992,
    developer: 'Capcom',
    imageUrl: 'https://i.redd.it/2rgdsgr7p3pc1.jpeg',
    launchConfig: { gameId: 'sf2ce' },
  },
  {
    id: '6',
    title: 'Guilty Gear Strive',
    platform: 'Steam',
    releaseYear: 2021,
    developer: 'Arc System Works',
    imageUrl: 'https://upload.wikimedia.org/wikipedia/en/7/7d/Guilty_Gear_Strive.jpg',
    launchConfig: { appId: '1384160' },
  },
  {
    id: '7',
    title: 'Tekken 8',
    platform: 'Steam',
    releaseYear: 2024,
    developer: 'Bandai Namco',
    imageUrl: 'https://cdn2.steamgriddb.com/grid/a9283100ad06971a29f5382f6ab25ea4.jpg',
    launchConfig: { appId: '1778820' },
  },
  {
    id: '8',
    title: 'Garou: Mark of the Wolves',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'SNK',
    imageUrl: 'https://m.media-amazon.com/images/M/MV5BNmM2OGEwMDItNzhjNy00MWJkLTljY2EtYTk3MjA5ZmM5ZWE1XkEyXkFqcGc@._V1_FMjpg_UX1000_.jpg',
    launchConfig: { gameId: 'garou' },
  },
  {
    id: '9',
    title: 'Samurai Shodown V Special',
    platform: 'Fightcade',
    releaseYear: 2004,
    developer: 'SNK',
    imageUrl: 'https://www.everyeye.it/public/immagini/14072017/samurai-shodown-v-special_notizia.jpg',
    launchConfig: { gameId: 'samsh5sp' },
  },
  {
    id: '10',
    title: 'Metal Slug 3',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'SNK',
    imageUrl: 'https://static.wikia.nocookie.net/metalslug/images/f/ff/IMG_20190220_164837.jpg/revision/latest?cb=20190220144903',
    launchConfig: { driverName: 'mslug3' },
  },
]

let nextId = 11

export async function fetchRecentlyPlayed(): Promise<Game[]> {
  return games
}

export async function fetchAllGames(query?: string): Promise<Game[]> {
  let result = [...games].sort((a, b) => a.title.localeCompare(b.title))
  if (query) {
    const q = query.toLowerCase()
    result = result.filter((g) => g.title.toLowerCase().includes(q))
  }
  return result
}

export async function fetchGame(id: string): Promise<Game | undefined> {
  return games.find((g) => g.id === id)
}

export async function createGame(game: Omit<Game, 'id'>): Promise<Game> {
  const newGame = { ...game, id: String(nextId++) }
  games.push(newGame)
  return newGame
}

export async function updateGame(id: string, data: Omit<Game, 'id'>): Promise<Game> {
  const idx = games.findIndex((g) => g.id === id)
  const updated = { ...data, id }
  games[idx] = updated
  return updated
}

export async function deleteGame(id: string): Promise<void> {
  const idx = games.findIndex((g) => g.id === id)
  if (idx !== -1) games.splice(idx, 1)
}

export async function openKeyboard(): Promise<void> {
  // TODO: replace with actual backend call that spawns steam-osk
}

const platformCatalog: Record<Platform, PlatformSearchResult[]> = {
  Steam: [
    { name: 'Street Fighter 6', platformId: '1364780' },
    { name: 'Street Fighter V', platformId: '310950' },
    { name: 'Street Fighter 30th Anniversary Collection', platformId: '586200' },
    { name: 'Guilty Gear Strive', platformId: '1384160' },
    { name: 'Guilty Gear Xrd REV 2', platformId: '520440' },
    { name: 'Tekken 8', platformId: '1778820' },
    { name: 'Tekken 7', platformId: '389730' },
    { name: 'Dragon Ball FighterZ', platformId: '678950' },
    { name: 'Mortal Kombat 1', platformId: '1971870' },
    { name: 'Mortal Kombat 11', platformId: '976310' },
    { name: 'The King of Fighters XV', platformId: '1498570' },
    { name: 'Under Night In-Birth II', platformId: '1797580' },
    { name: 'Blazblue Centralfiction', platformId: '586140' },
    { name: 'Granblue Fantasy Versus Rising', platformId: '2157560' },
    { name: 'DNF Duel', platformId: '1216060' },
    { name: 'Them\'s Fightin\' Herds', platformId: '574980' },
    { name: 'Rivals of Aether II', platformId: '2217000' },
    { name: 'Melty Blood Type Lumina', platformId: '1372280' },
    { name: 'Skullgirls 2nd Encore', platformId: '245170' },
    { name: 'Brawlhalla', platformId: '291550' },
  ],
  Fightcade: [
    { name: 'Street Fighter III: 3rd Strike', platformId: 'sfiii3n' },
    { name: 'Street Fighter II: Champion Edition', platformId: 'sf2ce' },
    { name: 'Street Fighter Alpha 3', platformId: 'sfa3' },
    { name: 'Super Street Fighter II Turbo', platformId: 'ssf2t' },
    { name: 'Hyper Street Fighter II', platformId: 'hsf2' },
    { name: 'King of Fighters 98', platformId: 'kof98' },
    { name: 'King of Fighters 2002', platformId: 'kof2002' },
    { name: 'King of Fighters 97', platformId: 'kof97' },
    { name: 'Garou: Mark of the Wolves', platformId: 'garou' },
    { name: 'Real Bout Fatal Fury 2', platformId: 'rbff2' },
    { name: 'Samurai Shodown V Special', platformId: 'samsh5sp' },
    { name: 'Samurai Shodown II', platformId: 'samsho2' },
    { name: 'JoJo\'s Bizarre Adventure: Heritage for the Future', platformId: 'jojoban' },
    { name: 'Vampire Savior', platformId: 'vsav' },
    { name: 'Marvel vs Capcom: Clash of Super Heroes', platformId: 'mvsc' },
    { name: 'X-Men vs Street Fighter', platformId: 'xmvsf' },
    { name: 'Capcom vs SNK 2', platformId: 'cvs2' },
    { name: 'Metal Slug 3', platformId: 'mslug3' },
    { name: 'Windjammers', platformId: 'wjammers' },
    { name: 'Puzzle Bobble', platformId: 'pbobblen' },
  ],
  MAME: [
    { name: 'Marvel vs Capcom 2', platformId: 'mvsc2' },
    { name: 'Street Fighter III: 3rd Strike', platformId: 'sfiii3' },
    { name: 'Metal Slug 3', platformId: 'mslug3' },
    { name: 'Metal Slug X', platformId: 'mslugx' },
    { name: 'Metal Slug', platformId: 'mslug' },
    { name: 'King of Fighters 2002', platformId: 'kof2002' },
    { name: 'Garou: Mark of the Wolves', platformId: 'garou' },
    { name: 'Pac-Man', platformId: 'pacman' },
    { name: 'Donkey Kong', platformId: 'dkong' },
    { name: 'Galaga', platformId: 'galaga' },
    { name: 'Space Invaders', platformId: 'invaders' },
    { name: 'Bubble Bobble', platformId: 'bublbobl' },
    { name: 'Mortal Kombat II', platformId: 'mk2' },
    { name: 'NBA Jam', platformId: 'nbajam' },
    { name: 'Teenage Mutant Ninja Turtles', platformId: 'tmnt' },
    { name: 'The Simpsons', platformId: 'simpsons' },
    { name: 'X-Men', platformId: 'xmen' },
    { name: 'Sunset Riders', platformId: 'ssriders' },
    { name: 'Cadillacs and Dinosaurs', platformId: 'dino' },
    { name: 'Dungeons & Dragons: Shadow over Mystara', platformId: 'ddsom' },
  ],
}

export async function searchPlatformGames(platform: Platform, query: string): Promise<PlatformSearchResult[]> {
  const q = query.toLowerCase()
  return platformCatalog[platform].filter((g) => g.name.toLowerCase().includes(q) || g.platformId.toLowerCase().includes(q))
}
