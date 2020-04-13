import { WarningOutlined } from '@ant-design/icons';
import classNames from 'classnames';
import { PREFIXCLS } from 'config/config';
import * as React from 'react';

interface TipsProps {
    tip: string
    icon?: string
    className?: string
    size?: 'large' | 'normal'
}

interface TipsStatus {
}

export default class Tips extends React.Component<TipsProps, TipsStatus> {
    render() {
        const { icon, tip, size, className } = this.props
        const prefixCls = `${PREFIXCLS}-tips`

        const classes = classNames(prefixCls, className, {
            [`${prefixCls}-large`]: size === 'large'
        })

        return (
            <div className={classes}>
                {icon && <span className={`${prefixCls}__icon`}><WarningOutlined /></span>}
                <span className={`${prefixCls}__content`}>{tip}</span>
            </div>
        )
    }
}