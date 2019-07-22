export const MENUS = [
  {
    title: 'Home',
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
        icon: 'share-alt',
        path: '/system',
      },
    ],
  }, {
    title: 'Setting',
    icon: 'setting',
    path: '/setting',
    children: [
      {
        title: 'Cluster',
        icon: 'share-alt',
        path: '/cluster',
      },
      {
        title: 'Logic Database',
        icon: 'database',
        path: '/logic/database',
      },
    ],
  },
]