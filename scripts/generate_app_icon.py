from __future__ import annotations

from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter


SIZE = 1024
ROOT = Path(__file__).resolve().parents[1]
BUILD_DIR = ROOT / "build"
WINDOWS_DIR = BUILD_DIR / "windows"
DOCS_ASSETS_DIR = ROOT / "docs" / "assets"


def lerp(a: int, b: int, t: float) -> int:
    return int(round(a + (b - a) * t))


def blend(c1: tuple[int, int, int], c2: tuple[int, int, int], t: float) -> tuple[int, int, int]:
    return tuple(lerp(c1[i], c2[i], t) for i in range(3))


def radial_background() -> Image.Image:
    center = (SIZE / 2, SIZE / 2)
    inner = (14, 20, 36)
    outer = (4, 6, 12)
    image = Image.new("RGBA", (SIZE, SIZE))
    pixels = image.load()

    for y in range(SIZE):
        for x in range(SIZE):
            dx = (x - center[0]) / SIZE
            dy = (y - center[1]) / SIZE
            distance = min(1.0, ((dx * dx + dy * dy) ** 0.5) * 1.85)
            base = blend(inner, outer, distance)

            tilt = max(0.0, min(1.0, (x * 0.65 + y * 0.35) / SIZE))
            accent = blend((12, 18, 30), (34, 16, 20), tilt)
            color = tuple(min(255, int(base[i] * 0.78 + accent[i] * 0.22)) for i in range(3))
            pixels[x, y] = (*color, 255)

    return image


def apply_rounded_mask(image: Image.Image, radius: int = 220) -> Image.Image:
    mask = Image.new("L", image.size, 0)
    ImageDraw.Draw(mask).rounded_rectangle((0, 0, SIZE - 1, SIZE - 1), radius=radius, fill=255)
    rounded = Image.new("RGBA", image.size, (0, 0, 0, 0))
    rounded.paste(image, mask=mask)
    return rounded


def bubble_layer(bounds: tuple[int, int, int, int], color_a: tuple[int, int, int], color_b: tuple[int, int, int], tail: list[tuple[int, int]], radius: int) -> Image.Image:
    layer = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    mask = Image.new("L", (SIZE, SIZE), 0)
    draw = ImageDraw.Draw(mask)
    draw.rounded_rectangle(bounds, radius=radius, fill=255)
    draw.polygon(tail, fill=255)

    gradient = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    gradient_pixels = gradient.load()
    left, top, right, bottom = bounds
    width = max(1, right - left)
    height = max(1, bottom - top)

    for y in range(top, bottom):
        for x in range(left, right):
            tx = (x - left) / width
            ty = (y - top) / height
            mixed = blend(color_a, color_b, min(1.0, max(0.0, tx * 0.72 + ty * 0.28)))
            gradient_pixels[x, y] = (*mixed, 255)

    ImageDraw.Draw(gradient).polygon(
        tail,
        fill=(*blend(color_a, color_b, 0.55), 255),
    )
    layer.paste(gradient, mask=mask)
    return layer


def add_shadow(base: Image.Image, shape_mask: Image.Image, color: tuple[int, int, int], blur: int, offset: tuple[int, int], alpha: int) -> None:
    shadow = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    solid = Image.new("RGBA", (SIZE, SIZE), (*color, alpha))
    shifted = Image.new("L", (SIZE, SIZE), 0)
    shifted.paste(shape_mask, offset)
    shadow.paste(solid, mask=shifted)
    shadow = shadow.filter(ImageFilter.GaussianBlur(blur))
    base.alpha_composite(shadow)


def make_shape_mask(draw_fn) -> Image.Image:
    mask = Image.new("L", (SIZE, SIZE), 0)
    draw = ImageDraw.Draw(mask)
    draw_fn(draw)
    return mask


def add_glow_ring(base: Image.Image) -> None:
    glow = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    draw.ellipse((188, 188, 836, 836), outline=(255, 88, 72, 95), width=28)
    draw.ellipse((230, 230, 794, 794), outline=(43, 205, 255, 65), width=14)
    draw.ellipse((262, 262, 762, 762), fill=(8, 12, 24, 175))
    glow = glow.filter(ImageFilter.GaussianBlur(6))
    base.alpha_composite(glow)


def add_center_spark(base: Image.Image) -> None:
    spark = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(spark)
    draw.polygon(
        [(492, 366), (560, 470), (668, 448), (584, 546), (632, 662), (512, 592), (398, 674), (438, 548), (338, 454), (460, 470)],
        fill=(255, 246, 235, 255),
    )
    draw.polygon(
        [(520, 400), (550, 470), (618, 462), (560, 524), (594, 600), (512, 558), (434, 614), (460, 534), (396, 468), (480, 476)],
        fill=(255, 146, 74, 255),
    )
    halo = spark.filter(ImageFilter.GaussianBlur(18))
    base.alpha_composite(halo)
    base.alpha_composite(spark)


def add_signal_ticks(base: Image.Image) -> None:
    layer = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    ticks = [
        ((162, 520, 226, 552), (255, 96, 73, 185)),
        ((814, 520, 878, 552), (46, 214, 255, 185)),
        ((492, 120, 532, 182), (255, 198, 82, 195)),
    ]
    for bounds, color in ticks:
        draw.rounded_rectangle(bounds, radius=22, fill=color)
    layer = layer.filter(ImageFilter.GaussianBlur(2))
    base.alpha_composite(layer)


def generate() -> Image.Image:
    base = apply_rounded_mask(radial_background())
    add_glow_ring(base)

    left_bounds = (140, 342, 438, 560)
    left_tail = [(250, 542), (164, 650), (338, 586)]
    right_bounds = (586, 332, 884, 550)
    right_tail = [(710, 534), (860, 622), (648, 578)]
    top_bounds = (356, 142, 668, 328)
    top_tail = [(502, 312), (452, 396), (566, 334)]

    left_mask = make_shape_mask(lambda draw: (draw.rounded_rectangle(left_bounds, radius=76, fill=255), draw.polygon(left_tail, fill=255)))
    right_mask = make_shape_mask(lambda draw: (draw.rounded_rectangle(right_bounds, radius=76, fill=255), draw.polygon(right_tail, fill=255)))
    top_mask = make_shape_mask(lambda draw: (draw.rounded_rectangle(top_bounds, radius=68, fill=255), draw.polygon(top_tail, fill=255)))

    add_shadow(base, top_mask, (255, 188, 60), blur=36, offset=(0, 18), alpha=96)
    add_shadow(base, left_mask, (255, 84, 72), blur=32, offset=(0, 18), alpha=96)
    add_shadow(base, right_mask, (45, 202, 255), blur=32, offset=(0, 18), alpha=96)

    base.alpha_composite(bubble_layer(top_bounds, (255, 211, 96), (255, 144, 54), top_tail, 68))
    base.alpha_composite(bubble_layer(left_bounds, (255, 115, 83), (216, 42, 73), left_tail, 76))
    base.alpha_composite(bubble_layer(right_bounds, (66, 223, 255), (47, 116, 255), right_tail, 76))

    inner = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    inner_draw = ImageDraw.Draw(inner)
    inner_draw.rounded_rectangle((286, 286, 738, 738), radius=144, outline=(255, 255, 255, 28), width=8)
    inner_draw.rounded_rectangle((314, 314, 710, 710), radius=126, outline=(255, 122, 80, 36), width=4)
    inner = inner.filter(ImageFilter.GaussianBlur(1))
    base.alpha_composite(inner)

    add_center_spark(base)
    add_signal_ticks(base)

    border = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    border_draw = ImageDraw.Draw(border)
    border_draw.rounded_rectangle((10, 10, SIZE - 10, SIZE - 10), radius=218, outline=(255, 255, 255, 54), width=4)
    border = border.filter(ImageFilter.GaussianBlur(0.3))
    base.alpha_composite(border)
    return base


def main() -> None:
    BUILD_DIR.mkdir(parents=True, exist_ok=True)
    WINDOWS_DIR.mkdir(parents=True, exist_ok=True)
    DOCS_ASSETS_DIR.mkdir(parents=True, exist_ok=True)

    icon = generate()
    appicon_path = BUILD_DIR / "appicon.png"
    docs_icon_path = DOCS_ASSETS_DIR / "appicon.png"
    windows_icon_path = WINDOWS_DIR / "icon.ico"

    icon.save(appicon_path, format="PNG")
    icon.save(docs_icon_path, format="PNG")
    icon.save(windows_icon_path, format="ICO", sizes=[(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (16, 16)])

    print(f"generated {appicon_path}")
    print(f"generated {docs_icon_path}")
    print(f"generated {windows_icon_path}")


if __name__ == "__main__":
    main()
