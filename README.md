# vrapi

mgl is used as the format for matrices, vectors and quarterions.

This causes two quirks. First is the matrices that the C API returns are in row major order. mgl (and opengl) expect this to be the transpose. (TODO should we return the transpose or just the matrix probaly just the transpose).

The second quirk is the quarterions used by the C API are using the JPL convention while mgl uses the hamilton (the one that is more standard across the scientific community).

???
The API will handle all the conversions behind the scenes so that whenever a matrix is returned (or a struct with a matrix) it will transpose the matrix. Whenever a quartenion is returned it will be in hamilton convention. Whenever the API expects a matrix or quartenion behind the scenes the API will trans from from hamilton to JPL. So you should be be able to work completly within mgls standards.  
    