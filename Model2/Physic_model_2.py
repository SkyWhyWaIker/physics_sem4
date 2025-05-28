import numpy as np
import matplotlib.pyplot as plt


class ApertureGenerator:
    @staticmethod
    def create_single_slit(x, y, width=0.1):
        return np.where(np.abs(x) <= width / 2, 1, 0)

    @staticmethod
    def create_rectangle(x, y, width=0.1, height=0.2):
        return np.where((np.abs(x) <= width / 2) & (np.abs(y) <= height / 2), 1, 0)

    @staticmethod
    def create_circular_aperture(x, y, diameter=1.0):
        r = np.sqrt(x ** 2 + y ** 2)
        return np.where(r <= diameter / 2, 1, 0)


class DiffractionCalculator:
    def __init__(self, L=0.5, k=600e-9, N=1000, D=0.5):
        self.L = L
        self.k = k
        self.N = N
        self.D = D
        self.dx = D / N

        x = np.linspace(-D / 2, D / 2, N)
        y = np.linspace(-D / 2, D / 2, N)
        self.x, self.y = np.meshgrid(x, y)

        self.apertures = {
            "single_slit": ApertureGenerator.create_single_slit(self.x, self.y, width=0.1),
            "rectangle": ApertureGenerator.create_rectangle(self.x, self.y, width=0.1, height=0.2),
            "circular": ApertureGenerator.create_circular_aperture(self.x, self.y, diameter=D)
        }

    def compute_fraunhofer_diffraction(self, amplitude):
        N = amplitude.shape[0]
        fft_result = np.fft.fftshift(np.fft.fft2(amplitude))
        print(f"Максимальное значение после FFT: {np.max(np.abs(fft_result))}")

        x_fraunhofer = np.fft.fftshift(np.fft.fftfreq(N, self.dx)) * self.k * self.L
        y_fraunhofer = np.fft.fftshift(np.fft.fftfreq(N, self.dx)) * self.k * self.L
        x_fraunhofer, y_fraunhofer = np.meshgrid(x_fraunhofer, y_fraunhofer)
        intensity = np.abs(fft_result) ** 2
        return intensity, x_fraunhofer, y_fraunhofer

    @staticmethod
    def plot_diffraction(amplitude, intensity, x_fraunhofer, y_fraunhofer, filename):
        plt.figure(figsize=(12, 6))

        plt.subplot(1, 2, 1)
        plt.imshow(amplitude, extent=(-0.5, 0.5, -0.5, 0.5), cmap='gray')
        plt.title('Амплитудное распределение')
        plt.colorbar()

        plt.subplot(1, 2, 2)
        plt.imshow(np.log1p(intensity),
                   extent=(x_fraunhofer.min(), x_fraunhofer.max(),
                           y_fraunhofer.min(), y_fraunhofer.max()),
                   cmap='inferno')
        plt.title('Дифракция Фраунгофера')
        plt.colorbar()

        plt.savefig(filename)
        plt.show()

    def process_all_apertures(self):
        for name, amplitude in self.apertures.items():
            intensity, x_f, y_f = self.compute_fraunhofer_diffraction(amplitude)
            self.plot_diffraction(amplitude, intensity, x_f, y_f, f"{name}.png")


if __name__ == "__main__":
    calculator = DiffractionCalculator()
    calculator.process_all_apertures()