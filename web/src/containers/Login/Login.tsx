import { Layout } from 'antd'
import * as React from 'react'

import LoginBackground from './LoginBackground'
import Footer from '../../components/Layout/Footer'

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
        <Content>
          Login
          <LoginBackground/>
        </Content>

        <Footer/>
      </Layout>
    )
  }
}