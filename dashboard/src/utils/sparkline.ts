export function generateSparklinePath(
  data: number[],
  width = 100,
  height = 24,
): string {
  if (data.length === 0) return "";

  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;

  const stepX = width / (data.length - 1 || 1);

  return data
    .map((value, index) => {
      const x = index * stepX;
      const y = height - ((value - min) / range) * height;
      return `${index === 0 ? "M" : "L"}${x.toFixed(2)},${y.toFixed(2)}`;
    })
    .join(" ");
}

export function generateSparklineArea(
  data: number[],
  width = 100,
  height = 24,
): string {
  if (data.length === 0) return "";

  const line = generateSparklinePath(data, width, height);
  return `${line} L${width},${height} L0,${height} Z`;
}

export function generateSparklinePoints(
  data: number[],
  width = 100,
  height = 24,
): { x: number; y: number }[] {
  if (data.length === 0) return [];

  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;
  const stepX = width / (data.length - 1 || 1);

  return data.map((value, index) => ({
    x: index * stepX,
    y: height - ((value - min) / range) * height,
  }));
}
