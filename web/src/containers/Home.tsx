/**
 * me
 * Created by jie.huang on 2018/7/18.
 */
import * as React from "react";
import {Tabs} from "antd";
import {observer} from "mobx-react";

const {LazyLog} = require("react-lazylog/build");
const R = require("ramda");
const TabPane = Tabs.TabPane;

interface HomeProps {
}
interface HomeState {
    zkDump?: string;
    threadDump?: string;
}

@observer
class Home extends React.Component<HomeProps, HomeState> {

    constructor(props: Readonly<HomeProps>) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <div>
                test
            </div>
        );
    }
}

export default Home;
