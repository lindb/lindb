import * as React from "react";
import {HashRouter as Router, Link, Redirect, Route, Switch} from "react-router-dom";
import {Breadcrumb, Card, Icon, Layout, Menu} from "antd";
import HomePage from "./containers/Home";
import {autobind} from "core-decorators";
import "./App.css";
import {MENUS} from "./config";

const {Sider, Content, Footer, Header} = Layout;
const MenuItem = Menu.Item;

interface AppProps {
}

interface AppState {
    collapsed: boolean;
}
class App extends React.Component<AppProps, AppState> {

    constructor(props: Readonly<AppProps>) {
        super(props);
        this.state = {collapsed: false};
    }

    @autobind
    toggle() {
        this.setState({
            collapsed: !this.state.collapsed,
        });
    }

    renderMenu(parent: string, menus: Array<any>) {
        return menus.map(menu => {
            if (menu.children) {
                return (
                    <Menu.SubMenu
                        key={menu.path}
                        title={<span><Icon type={menu.icon}/><span>{menu.text}</span></span>}
                    >
                        {this.renderMenu(menu.path, menu.children)}
                    </Menu.SubMenu>
                );
            } else {
                const path = parent ? parent + menu.path : menu.path;
                return (
                    <MenuItem key={path}>
                        <Link to={path}>
                            <Icon type={menu.icon}/>
                            <span>{menu.text}</span>
                        </Link>
                    </MenuItem>
                );
            }
        });
    }

    renderBreadcrumb(parent: string, menus: Array<any>, path: string, breadcrumb: Array<any>) {
        return menus.map(menu => {
            const uri = parent ? parent + menu.path : menu.path;
            if (path.startsWith(uri)) {
                breadcrumb.push(<Breadcrumb.Item key={uri}>{menu.text}</Breadcrumb.Item>);
                if (menu.children) {
                    this.renderBreadcrumb(menu.path, menu.children, path, breadcrumb);
                }
            }
        });
    }

    public render() {
        const {collapsed} = this.state;
        const breadcrumb = [];
        const {location: {hash}} = window;
        const path = hash.replace("#", "");
        this.renderBreadcrumb(null, MENUS, path, breadcrumb);
        return (
            <Layout style={{height: "100vh"}}>
                <Sider style={{overflowY: "scroll", zIndex: 99}} collapsed={collapsed}>
                    <Menu
                        theme="light"
                        mode="inline"
                        defaultOpenKeys={["monitoring", "/setting"]}
                        defaultSelectedKeys={[path]}
                        style={{height: "100vh"}}
                    >
                        <MenuItem style={{borderColor: "#f8f8f8"}} disabled={true}>
                            <img src="images/LinDB.png" alt="logo" style={{height: 38}}/>
                        </MenuItem>
                        {this.renderMenu(null, MENUS)}
                    </Menu>
                </Sider>
                <Layout>
                    <Header className="global-header-index-header">
                        <Icon
                            className="collapse"
                            type={collapsed ? "menu-unfold" : "menu-fold"}
                            onClick={this.toggle}
                        />
                    </Header>
                    <Content style={{margin: 6, marginBottom: 70, paddingTop: 50}}>
                        <Card style={{marginBottom: 6}}>
                            <Breadcrumb>
                                {breadcrumb.map(item => item)}
                            </Breadcrumb>
                        </Card>
                        <Router>
                            <Switch>
                                <Route exact={true} path="/" component={HomePage}/>
                                <Route exact={true} path="/search" component={HomePage}/>
                                <Route exact={false} path="/monitoring" component={HomePage}/>
                                <Route exact={false} path="/setting" component={HomePage}/>
                                <Redirect to="/"/>
                            </Switch>
                        </Router>
                    </Content>
                    <Footer
                        style={{
                            textAlign: "center",
                            padding: 8,
                            backgroundColor: "#f5f7fa",
                            position: "fixed",
                            bottom: 0,
                            left: 0,
                            borderTop: "1px dotted #2aabd2",
                            right: 0
                        }}
                    >
                        LinDB © 2018 Created by Framework Team ele.me
                    </Footer>
                </Layout>
            </Layout>
        );

    }
}

export default App;
