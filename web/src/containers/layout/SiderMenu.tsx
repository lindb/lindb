import { DatabaseOutlined, HomeOutlined, SearchOutlined, DeploymentUnitOutlined, HddOutlined, DesktopOutlined, SettingOutlined, ShareAltOutlined } from '@ant-design/icons'
import { Layout, Menu } from 'antd'
import Logo from 'assets/images/logo_title_dark.png'
import { BizIcon } from 'components/basic/BizIcon'
import { MENUS } from 'config/menu'
import { autobind } from 'core-decorators'
import { BreadcrumbStatus } from 'model/Breadcrumb'
import * as React from 'react'
import { Link } from 'react-router-dom'
import StoreManager from 'store/StoreManager'

const { Sider } = Layout

interface SiderMenuProps {
}

interface SiderMenuStatus {
}

export default class SiderMenu extends React.Component<SiderMenuProps, SiderMenuStatus> {
    breadcrumbStore: any
    flatMenu: Array<any>
    currentBreadcrumbPath: Array<any>
    iconMap: Map<String, any>;

    constructor(props: SiderMenuProps) {
        super(props)
        this.breadcrumbStore = StoreManager.BreadcrumbStore
        this.flatMenu = this.getFlatMenu()
        this.currentBreadcrumbPath = []
        this.iconMap = new Map()
        this.iconMap.set("home", <HomeOutlined />)
        this.iconMap.set("search", <SearchOutlined />)
        this.iconMap.set("share-alt", <ShareAltOutlined />)
        this.iconMap.set("database", <DatabaseOutlined />)
        this.iconMap.set("setting", <SettingOutlined />)
        this.iconMap.set("system", <DesktopOutlined />)
        this.iconMap.set("broker", <DeploymentUnitOutlined />)
        this.iconMap.set("storage", <HddOutlined />)
        this.iconMap.set("runtime", <BizIcon type={"icongo"} style={{ fontSize: '16px' }} />)
    }

    renderMenu(menus: Array<any>, parentPath?: string) {
        const IconTitle = (icon: string, title: string) => (<span>{this.iconMap.get(icon)}{title}</span>)
        return menus.map(menu => {
            const path = parentPath ? parentPath + menu.path : menu.path

            return menu.children
                ? (
                    <Menu.SubMenu key={menu.path} title={IconTitle(menu.icon, menu.title)}>
                        {this.renderMenu(menu.children, menu.path)}
                    </Menu.SubMenu>
                )
                : (
                    <Menu.Item key={path}>
                        <Link to={path}>{IconTitle(menu.icon, menu.title)}</Link>
                    </Menu.Item>
                )
        })
    }

    @autobind
    handleMenuClick(e: any) {
        /* e.keyPath为点击菜单所在的path以及parent path组成的数组，parent path在数组尾部*/
        this.currentBreadcrumbPath = e.keyPath.reverse()
        this.setBreadcrumbs()
    }

    /**
     * @description update the data in breadcrumbStore based on the current breadcrumb routing data
     */
    @autobind
    setBreadcrumbs() {
        let breadcrumbs: Array<BreadcrumbStatus> = []
        this.currentBreadcrumbPath.forEach((path: string) => (
            this.flatMenu.forEach(item => {
                if (path === item.path) {
                    breadcrumbs.push(item)
                }
            })
        ))
        this.breadcrumbStore.setBreadcrumbs(breadcrumbs)
    }

    /**
     * @description Flatten the parent-child structure of the MENU to fit Breadcrumb structure
     * @return flattening menu
     */
    @autobind
    getFlatMenu() {
        let flatMenu: Array<BreadcrumbStatus> = []
        MENUS.forEach(item => {
            if (item.children) {
                flatMenu.push({ path: item.path, label: item.title })
                item.children.forEach(child => {
                    flatMenu.push({ path: item.path + child.path, label: child.title });
                })
            } else {
                flatMenu.push({ path: item.path, label: item.title })
            }
        })
        return flatMenu
    }
    /**
     * @description initialize the breadcrumb data by getting the page routing
     */
    @autobind
    initBreadcrumb() {
        const path = this.getPath()
        let pathArr = path.split('/').slice(1)
        let result = ''
        let breadcrumbPathArray = pathArr.map((e, i) => {
            result += '/' + e
            return result
        })

        this.currentBreadcrumbPath = breadcrumbPathArray
        this.setBreadcrumbs()
    }

    /**
     * @description get current hash path and get rid of '#'
     * @example '#/index' => '/index'
     */
    @autobind
    getPath() {
        const { location: { hash } } = window;
        const path = hash.replace('#', '')
        return path
    }

    componentDidMount(): void {
        this.initBreadcrumb()
    }

    render() {
        const path = this.getPath()
        return (
            <Sider className="lindb-sider" collapsible={true} trigger={null}>
                {/* Logo */}
                <div className="lindb-sider__logo">
                    <Link to="/"><img src={Logo} alt="LinDB" /></Link>
                </div>

                {/* Menu */}
                <Menu
                    mode="inline"
                    theme="dark"
                    className="lindb-sider__menu"
                    defaultOpenKeys={['/monitoring', '/metadata']}
                    selectedKeys={[path]}
                    onClick={this.handleMenuClick}
                >
                    {this.renderMenu(MENUS)}
                </Menu>
            </Sider>
        )
    }
}