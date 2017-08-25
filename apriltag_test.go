package apriltag

import (
  "image"
  "image/color"
  "image/png"
  "os"
  "testing"
)


func TestCreateClose(t *testing.T){
  detector := New()
  err := detector.Close()
  if err != nil {
    t.Error(err)
  }
}

func TestFormatTransform(t *testing.T){
  f, err := os.Open("testtags.png")
  if err != nil {
    t.Fatal(err)
  }
  defer f.Close()

  img, err := png.Decode(f)
  if err != nil {
    t.Fatal(err)
  }

  cImg := grayToC(ImgToGrayscale(img))
  if cImg == nil {
    t.Error("Expected non-nil image")
  }
}


func TestFind(t *testing.T){
  detector := New()
  defer detector.Close()
  f, err := os.Open("testtags.png")
  if err != nil {
    t.Fatal(err)
  }
  defer f.Close()

  img, err := png.Decode(f)
  if err != nil {
    t.Fatal(err)
  }

  findings := detector.Find(ImgToGrayscale(img))
  t.Log(findings)

  center := color.RGBA{R: 255, G: 0, B: 0, A: 255}
  corner := color.RGBA{R: 0, G: 255, B: 0, A: 255}
  DrawFindings(img.(*image.RGBA), findings, center, corner)
  fOut, err := os.Create("test_output.png")
  if err != nil {
    t.Fatal(err)
  }
  defer fOut.Close()

  if err = png.Encode(fOut, img); err != nil {
    t.Fatal(err)
  }
}
