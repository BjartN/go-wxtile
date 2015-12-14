package tile


func BilinearInterpolation (ll *PointValue , ul *PointValue , lr *PointValue , ur *PointValue , p *Point) float64 {
	interpolatedValue1:= straightLineInterpolation(ll.X, float64(ll.Value), lr.X,  float64(lr.Value), p.X)
	interpolatedValue2 := straightLineInterpolation(ul.X,  float64(ul.Value), ur.X,  float64(ur.Value), p.X)

	return straightLineInterpolation(ll.Y, interpolatedValue1, ul.Y, interpolatedValue2, p.Y)
}

func straightLineInterpolation( x1 float64,  v1 float64,  x2 float64,  v2 float64,  x float64) float64{

	if(x1 == x){
		return v1
	}

	if(x2==x){
		return v2
	}

	if x > x2 {
		if (x - x2) < .00005{
			x = x2;
		}else  {
			return (v1 + v2) / 2;
		}
	}

	if x < x1{
		if ((x1 - x) < .00005){
			x = x1;
		}else {
			return (v1+v2)/2;
		}
	}

	dx := x2 - x1;


	return ((x2 - x) / dx) * v1 + ((x - x1) / dx) * v2;
}