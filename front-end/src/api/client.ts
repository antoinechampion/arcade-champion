import type { Game } from './types'

const mockRecentlyPlayed: Game[] = [
  {
    id: '1',
    title: 'Street Fighter III: 3rd Strike',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'Capcom',
    imageUrl: 'https://i.pinimg.com/736x/a8/db/46/a8db46f121860572350471b0e4405c32.jpg',
    bannerUrl: 'https://d1lss44hh2trtw.cloudfront.net/resize?type=webp&url=https%3A%2F%2Fshacknews-www.s3.amazonaws.com%2Fassets%2Farticle%2F2024%2F07%2F29%2F3rd-strike-header_feature.jpg&width=1032&sign=XEMEx9JTW7qcjq1YKxMlyIOL40hjUM0eBKwjqW6KlMU',
  },
  {
    id: '2',
    title: 'Marvel VS Capcom 2',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'Capcom',
    imageUrl: 'https://wiki.supercombo.gg/images/thumb/2/29/MVSC2_Cover_Art.jpg/300px-MVSC2_Cover_Art.jpg',
  },
  {
    id: '3',
    title: 'King of Fighters 98',
    platform: 'Fightcade',
    releaseYear: 1998,
    developer: 'SNK',
    imageUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1y8h.jpg',
  },
  {
    id: '4',
    title: 'Street Fighter 6',
    platform: 'Steam',
    releaseYear: 2023,
    developer: 'Capcom',
    imageUrl: 'https://upload.wikimedia.org/wikipedia/en/thumb/9/94/Street_Fighter_6_box_art.jpg/250px-Street_Fighter_6_box_art.jpg',
  },
  {
    id: '5',
    title: 'Street Fighter 2: Champion Edition',
    platform: 'Fightcade',
    releaseYear: 1992,
    developer: 'Capcom',
    imageUrl: 'https://i.redd.it/2rgdsgr7p3pc1.jpeg',
  },
  {
    id: '6',
    title: 'Guilty Gear Strive',
    platform: 'Steam',
    releaseYear: 2021,
    developer: 'Arc System Works',
    imageUrl: 'https://upload.wikimedia.org/wikipedia/en/7/7d/Guilty_Gear_Strive.jpg',
  },
  {
    id: '7',
    title: 'Tekken 8',
    platform: 'Steam',
    releaseYear: 2024,
    developer: 'Bandai Namco',
    imageUrl: 'https://cdn2.steamgriddb.com/grid/a9283100ad06971a29f5382f6ab25ea4.jpg',
  },
  {
    id: '8',
    title: 'Garou: Mark of the Wolves',
    platform: 'Fightcade',
    releaseYear: 1999,
    developer: 'SNK',
    imageUrl: 'https://m.media-amazon.com/images/M/MV5BNmM2OGEwMDItNzhjNy00MWJkLTljY2EtYTk3MjA5ZmM5ZWE1XkEyXkFqcGc@._V1_FMjpg_UX1000_.jpg',
  },
  {
    id: '9',
    title: 'Samurai Shodown V Special',
    platform: 'Fightcade',
    releaseYear: 2004,
    developer: 'SNK',
    imageUrl: 'https://www.everyeye.it/public/immagini/14072017/samurai-shodown-v-special_notizia.jpg',
  },
  {
    id: '10',
    title: 'Metal Slug 3',
    platform: 'MAME',
    releaseYear: 2000,
    developer: 'SNK',
    imageUrl: 'https://static.wikia.nocookie.net/metalslug/images/f/ff/IMG_20190220_164837.jpg/revision/latest?cb=20190220144903',
  },
]

export async function fetchRecentlyPlayed(): Promise<Game[]> {
  // TODO: replace with actual backend call
  return mockRecentlyPlayed
}

export async function fetchAllGames(query?: string): Promise<Game[]> {
  // TODO: replace with actual backend call
  let games = [...mockRecentlyPlayed].sort((a, b) => a.title.localeCompare(b.title))
  if (query) {
    const q = query.toLowerCase()
    games = games.filter((g) => g.title.toLowerCase().includes(q))
  }
  return games
}
