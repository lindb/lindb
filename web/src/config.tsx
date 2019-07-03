export const MENUS = [
    {
        text: "Home",
        icon: "home",
        path: "/home"
    }, {
        text: "Search",
        icon: "search",
        path: "/search"
    }, {
        text: "Setting",
        icon: "setting",
        path: "/setting",
        children: [
            {
                text: "Cluster",
                icon: "share-alt",
                path: "/cluster"
            },
            {
                text: "Logic Database",
                icon: "database",
                path: "/logic/database"
            }
        ]
    }
];