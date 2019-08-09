import { Button, Icon, Input, Layout } from 'antd'
import * as React from 'react'

import LoginBackground from './LoginBackground'
import Footer from '../../components/Layout/Footer'

const Logo = require('../../assets/images/logo_title.png')

const { Content } = Layout

interface LoginProps {
}

interface LoginStatus {
}

export default class Login extends React.Component<LoginProps, LoginStatus> {
  constructor(props: LoginProps) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <Layout className="lindb-login">
        <LoginBackground/>

        <Content>

          <div className="lindb-login__content">
            <img className="lindb-login__content__logo" src={Logo} alt=""/>
            <Input
              className="lindb-login__content__input"
              size="large"
              prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }}/>}
              placeholder="Enter Username"
            />

            <Input.Password
              className="lindb-login__content__input"
              size="large"
              prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }}/>}
              placeholder="Enter Password"
            />

            <Button
              className="lindb-login__content__btn"
              type="primary"
              shape="round"
              size="large"
              block={true}
            >
              Login
            </Button>
          </div>
        </Content>

        <Footer/>
      </Layout>
    )
  }
}