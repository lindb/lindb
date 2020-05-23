import { Chart } from 'model/Chart'
import { ChartDatasets, ResultSet } from 'model/Metric'
import { getChartColor, toRGBA } from 'utils/Util'

const ChartJS = require('chart.js')
const moment = require('moment')
const Color = ChartJS.helpers.color
/**
 * Generate Line Chart data and options
 * @param {ResultSet} resultSet
 * @param {UnitEnum} unit Current chart Y-axes unit
 */
export function LineChart(resultSet: ResultSet | null, chart: Chart, selectedSeries: Map<string, string>) {
  if (!resultSet) {
    return {}
  }

  const { series, startTime, endTime, interval } = resultSet

  if (!series || series.length === 0) {
    return {}
  }

  const datasets: ChartDatasets[] = []
  let colorIdx = 0
  let leftMax = 0;

  series.forEach(item => {
    const { tags, fields } = item

    if (!fields) {
      return
    }

    const groupName = getGroupByTags(tags)

    for (let key of Object.keys(fields)) {
      const bgColor = getChartColor(colorIdx++)

      const fill = chart.type === 'area'
      const borderColor = bgColor
      const backgroundColor = chart.type === 'area' ? Color(bgColor).alpha(0.25).rgbString() : bgColor;
      const label = groupName ? groupName : key
      const pointBackgroundColor = toRGBA(bgColor, 0.25)

      let data: any = []
      const points: { [timestamp: string]: number } = fields[key]
      let i = 0;
      let timestamp = startTime! + i * interval!
      for (; timestamp <= endTime!;) {
        const value = points[`${timestamp}`];
        const v = value ? Math.floor(value * 1000) / 1000 : 0
        if (leftMax < v) {
          leftMax = v
        }
        data.push(v)
        i++
        timestamp = startTime! + i * interval!
      }

      let hidden = false
      if (selectedSeries) {
        hidden = !selectedSeries.has(label)
      }
      datasets.push({ label, data, fill, backgroundColor, borderColor, pointBackgroundColor, hidden })
    }
  })
  const labels = [];
  const start = new Date(startTime!);
  const end = new Date(endTime!);
  const showTimeLabel = start.getDate() !== end.getDate() || start.getMonth() !== end.getMonth() || start.getFullYear() !== end.getFullYear();
  const range = endTime! - startTime!
  let i = 0
  let timestamp = startTime! + i * interval!
  const times = []
  for (; timestamp <= endTime!;) {
    const dateTime = moment(timestamp);
    if (showTimeLabel) {
      labels.push(dateTime.format("MM/DD HH:mm"));
    } else if (range > 5 * 60 * 1000) {
      labels.push(dateTime.format("HH:mm"));
    } else {
      labels.push(dateTime.format("HH:mm:ss"));
    }
    times.push(timestamp)
    i++
    timestamp = startTime! + i * interval!
  }
  return { labels, datasets, interval, times, leftMax }
}

function getGroupByTags(tags: any) {
  if (!tags) {
    return ""
  }
  const tagKeys = Object.keys(tags)
  if (tagKeys.length === 1) {
    return tags[tagKeys[0]]
  }
  const result = []
  for (let key of tagKeys) {
    result.push(key + ":" + tags[key])
  }
  return result.join(",")
}