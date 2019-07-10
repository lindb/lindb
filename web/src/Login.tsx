/**
 * Created by jie.huang on 2018/7/26.
 */
import * as React from "react";
import {Alert, Button, Card, Col, Form, Icon, Input, Layout, Row} from "antd";
import {autobind} from "core-decorators";

const {Sider, Content, Footer} = Layout;
const FormItem = Form.Item;

interface LoginProps {
    form?: any;
}

interface LoginState {
    loading?: boolean;
    error?: string;
}

class LoginForm extends React.Component<any, LoginState> {

    constructor(props: Readonly<LoginProps>) {
        super(props);
        this.state = {loading: false};
    }

    @autobind
    handleSubmit(e: any) {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (err) {
                return;
            }
            this.setState({loading: true, error: null});
            // LoginService.user(values).then(value => {
            //     if (value && value.response && value.response.status == 400) {
            //         this.setState({loading: false, error: value.response.data});
            //     } else {
            //         const url = window.location.hash;
            //         let forward = "/";
            //         const user = "#/user?from=#";
            //         if (url && url.indexOf(user) >= 0) {
            //             forward = url.substring(url.indexOf(user) + user.length);
            //         }
            //         window.location.hash = forward;
            //     }
            // }).catch(error => {
            //     this.setState({loading: false, error: JSON.stringify(error)});
            // });
        });
    }

    render() {
        const {getFieldDecorator} = this.props.form;
        const {loading, error} = this.state;
        return (
            <Layout>
                <Content style={{height: "100vh"}}>
                    <Row style={{paddingTop: 100}} onSubmit={this.handleSubmit}>
                        <Col span={8} offset={8}>
                            <Card>
                                <div style={{textAlign: "center", paddingBottom: 20}}>
                                    <img src="images/LinDB.png"/>
                                </div>
                                <Form className="login-form">
                                    {error && (
                                        <FormItem style={{textAlign: "center"}}>
                                            <Alert description={error} type="error" message="Login Error"/>
                                        </FormItem>
                                    )}
                                    <FormItem>
                                        {getFieldDecorator("username", {
                                            rules: [{required: true, message: "Please input your username!"}],
                                        })(
                                            <Input
                                                size="large"
                                                prefix={<Icon type="user" style={{color: "rgba(0,0,0,.25)"}}/>}
                                                placeholder="Input username"
                                            />
                                        )}
                                    </FormItem>
                                    <FormItem>
                                        {getFieldDecorator("password", {
                                            rules: [{required: true, message: "Please input your Password!"}],
                                        })(
                                            <Input
                                                size="large"
                                                prefix={<Icon type="lock" style={{color: "rgba(0,0,0,.25)"}}/>}
                                                type="password"
                                                placeholder="Input password"
                                            />
                                        )}
                                    </FormItem>
                                    <FormItem style={{textAlign: "center"}}>
                                        <Button
                                            type="primary"
                                            size="large"
                                            icon="login"
                                            htmlType="submit"
                                            loading={loading}
                                            className="login-form-button"
                                        >
                                            Log in
                                        </Button>
                                    </FormItem>
                                </Form>
                            </Card>
                        </Col>
                    </Row>
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
        );
    }
}

const Login = Form.create<LoginProps>({})(LoginForm);
export default Login;
