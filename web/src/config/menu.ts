export const MENUS = [
  {
    title: 'Overview',
    icon: 'home',
    path: '/',
  }, {
    title: 'Search',
    icon: 'search',
    path: '/search',
  }, {
    title: 'Monitoring',
    icon: 'setting',
    path: '/monitoring',
    children: [
      {
        title: 'Runtime',
        icon: 'runtime',
        path: '/Runtime',
      },
    ],
  }, {
    title: 'Admin',
    icon: 'setting',
    path: '/admin',
    children: [
      {
        title: 'Storage',
        icon: 'share-alt',
        path: '/storage',
      },
      {
        title: 'Database',
        icon: 'database',
        path: '/database',
      },
    ],
  },
]