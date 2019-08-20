import * as React from 'react'

interface LoginBackgroundProps {
}

interface LoginBackgroundStatus {
}

export default class LoginBackground extends React.Component<LoginBackgroundProps, LoginBackgroundStatus> {
  canvas: React.RefObject<HTMLCanvasElement>

  NUM_PARTICLES = 300
  PARTICLE_SIZE = 0.5 // View heights
  SPEED = 10000 // Milliseconds

  particles: any[] = []

  constructor(props: LoginBackgroundProps) {
    super(props)
    this.canvas = React.createRef()
  }

  componentDidMount(): void {
    this.startAnimation()
  }

  // modified version of random-normal
  normalPool(o: any) {
    let r = 0
    do {
      const a = Math.round(this.randomNormal({ mean: o.mean, dev: o.dev }))
      if (a < o.pool.length && a >= 0) {
        return o.pool[a]
      }
      r ++
    } while (r < 100)
  }

  randomNormal(o: any) {
    // eslint-disable-next-line
    if (o = Object.assign({ mean: 0, dev: 1, pool: [] }, o), Array.isArray(o.pool) && o.pool.length > 0) {
      return this.normalPool(o)
    }
    let r, a, n, e, l = o.mean, t = o.dev
    do {
      r = (a = 2 * Math.random() - 1) * a + (n = 2 * Math.random() - 1) * n
    } while (r >= 1)
    // eslint-disable-next-line
    return e = a * Math.sqrt(-2 * Math.log(r) / r), t * e + l
  }

  rand(low: number, high: number) {
    return Math.random() * (high - low) + low
  }

  createParticle(): object {
    const colour = {
      r: 1,
      g: this.randomNormal({ mean: 150, dev: 20 }),
      b: this.randomNormal({ mean: 210, dev: 20 }),
      a: this.rand(0, 1)
    }

    return {
      x: -2,
      y: -2,
      diameter: Math.max(0, this.randomNormal({ mean: this.PARTICLE_SIZE, dev: this.PARTICLE_SIZE / 2 })),
      duration: this.randomNormal({ mean: this.SPEED, dev: this.SPEED * 0.1 }),
      amplitude: this.randomNormal({ mean: 16, dev: 2 }),
      offsetY: this.randomNormal({ mean: 0, dev: 10 }),
      arc: Math.PI * 2,
      startTime: performance.now() - this.rand(0, this.SPEED),
      colour: `rgba(${colour.r}, ${colour.g}, ${colour.b}, ${colour.a})`,
    }
  }

  moveParticle(particle: any, canvas: any, time: number) {
    const progress = ((time - particle.startTime) % particle.duration) / particle.duration
    return {
      ...particle,
      x: progress,
      y: ((Math.sin(progress * particle.arc) * particle.amplitude) + particle.offsetY),
    }
  }

  drawParticle(particle: any, canvas: any, ctx: any) {
    const vh = canvas.height / 100

    ctx.fillStyle = particle.colour
    ctx.beginPath()
    ctx.ellipse(
      particle.x * canvas.width,
      particle.y * vh + (canvas.height / 2),
      particle.diameter * vh,
      particle.diameter * vh,
      0,
      0,
      2 * Math.PI,
    )
    ctx.fill()
  }

  draw(time: number, canvas: any, ctx: any) {
    // Move particles
    this.particles.forEach((particle, index) => {
      this.particles[ index ] = this.moveParticle(particle, canvas, time)
    })

    // Clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // Draw the particles
    this.particles.forEach((particle) => {
      this.drawParticle(particle, canvas, ctx)
    })

    // Schedule next frame
    requestAnimationFrame((t: number) => this.draw(t, canvas, ctx))
  }

  initializeCanvas() {
    const canvas = this.canvas.current
    if (!canvas) {
      return [null, null]
    }

    canvas.width = canvas.offsetWidth * window.devicePixelRatio
    canvas.height = canvas.offsetHeight * window.devicePixelRatio

    let ctx = canvas.getContext('2d')

    window.addEventListener('resize', () => {
      canvas.width = canvas.offsetWidth * window.devicePixelRatio
      canvas.height = canvas.offsetHeight * window.devicePixelRatio
      ctx = canvas.getContext('2d')
    })

    return [ canvas, ctx ]
  }

  startAnimation() {
    const [ canvas, ctx ] = this.initializeCanvas()

    // Create a bunch of particles
    for (let i = 0; i < this.NUM_PARTICLES; i++) {
      this.particles.push(this.createParticle())
    }

    requestAnimationFrame((time) => this.draw(time, canvas, ctx))
  }

  render() {
    return (
      <div className="lindb-login__background">
        <canvas ref={this.canvas}/>
      </div>
    )
  }
}