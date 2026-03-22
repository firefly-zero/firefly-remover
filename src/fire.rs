use alloc::vec::Vec;
use firefly_rust::*;

pub struct Fire {
    force: f32,
    forced: f32,
    particles: Vec<Particle>,
}

impl Fire {
    pub fn new() -> Self {
        let mut particles = Vec::with_capacity(800);
        for _ in 0..particles.capacity() {
            let mut particle = Particle::new();
            particle.r = rand_f32() * 8.;
            particles.push(particle);
        }
        Self {
            force: 3.,
            forced: -1.,
            particles,
        }
    }

    pub fn update(&mut self) {
        self.force += 0.02 * self.forced;
        if self.force > 4. || self.force < -4. {
            self.forced *= -1.;
        }
        for particle in &mut self.particles {
            particle.update(self.force);
            if particle.r < 0. {
                *particle = Particle::new();
            }
        }
    }

    pub fn draw(&mut self) {
        for particle in &mut self.particles {
            particle.draw();
        }
    }
}

struct Particle {
    x: f32,
    y: f32,
    s: f32,
    r: f32,
}

impl Particle {
    fn new() -> Self {
        Self {
            x: rand_f32() * 240.,
            y: 150. + rand_f32() * 8.,
            s: 0.3 + rand_f32(),
            r: 12.,
        }
    }

    fn update(&mut self, force: f32) {
        self.x += force / 10.;
        self.y -= self.s;
        self.r -= self.s / 12.;
    }

    fn draw(&mut self) {
        let ax = (self.x - 120.).abs();
        let mut c = Color::White;

        if ax > 20. {
            c = Color::Yellow;
        }
        if ax > 30. {
            c = Color::Orange;
        }
        if self.r < 10. {
            c = Color::Yellow;
        }
        if self.r < 7. {
            c = Color::Orange;
        }
        if ax > 25. && c == Color::White {
            c = Color::Orange;
        }

        self.r -= ax / 480.;
        circle(self.x, self.y, self.r, c);
    }
}

fn circle(x: f32, y: f32, d: f32, c: Color) {
    let p = Point::new(x as i32, y as i32);
    draw_circle(p, d as i32, Style::solid(c));
}

fn rand_f32() -> f32 {
    get_random() as f32 / u32::MAX as f32
}
