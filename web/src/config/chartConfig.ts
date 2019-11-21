import { UnitEnum } from '../model/Metric'
import { DataFormatter } from '../utils/DataFormatter'

const ChartJS = require('chart.js')

// Global Chart Config
ChartJS.defaults.global.legend.display = false
ChartJS.defaults.global.tooltips.enabled = false
ChartJS.defaults.global.hover.mode = 'index'
ChartJS.defaults.global.hover.intersect = false
ChartJS.defaults.global.elements.point.radius = 0
ChartJS.defaults.global.elements.point.hoverRadius = 4
ChartJS.defaults.global.elements.line.borderWidth = 1

/**
 * Get ChartJS Options
 * @param {UnitEnum} unit Current chart Y-axes unit
 * @return chart config
 */
export function getOptions(unit?: UnitEnum) {
  if (!unit) {
    return
  }

  // Generate Options
  const color = { border: '#E3E3E3' }

  const scales = {
    yAxes: [{
      ticks: {
        fontColor: color.border,
        fontSize: 12,
        mirror: true, // tick in chart
        maxTicksLimit: 6,
        beginAtZero: true,
        tickMarkLength: 0,
        callback: function (value: number) {
          return DataFormatter.formatter(value, unit)
        },
      },
      gridLines: {
        lineWidth: 0.5,
        drawTicks: false,
        color: color.border,
        zeroLineColor: color.border,
      },
    }],
    xAxes: [{
      type: 'time',
      distribution: 'series',
      ticks: {
        fontColor: color.border,
        fontSize: 12,
      },
      gridLines: {
        lineWidth: 0.5,
        drawTicks: false,
        color: color.border,
        zeroLineColor: color.border,
      },
      time: {
        displayFormats: {
          second: 'HH:mm:ss',
          minute: 'HH:mm',
          hour: 'Do.H',
        },
      },
    }],
  }

  const elements = {
    line: {
      tension: 0, // disables bezier curves
    },
  }

  // disable animation config
  const disableAnimation = {
    animation: {
      duration: 0, // general animation time
    },
    hover: {
      animationDuration: 0, // duration of animations when hovering an item
    },
    responsiveAnimationDuration: 0, // animation duration after a resize
  }

  return {
    scales,
    elements,
    ...disableAnimation,
    // In order for Chart.js to obey the custom size you need to set maintainAspectRatio to false, example:
    responsive: true,
    maintainAspectRatio: false,
  }
}