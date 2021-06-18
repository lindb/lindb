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
        title: 'System',
        icon: 'system',
        path: '/system',
      },
      {
        title: 'Runtime',
        icon: 'runtime',
        path: '/runtime',
      },
      {
        title: 'Broker',
        icon: 'broker',
        path: '/broker',
      },
      {
        title: 'Storage',
        icon: 'storage',
        path: '/storage',
      },
    ],
  }, {
    title: 'Metadata',
    icon: 'setting',
    path: '/metadata',
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