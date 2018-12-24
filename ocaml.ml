let move l x = 
    let left x = x - 1 in
    let right x = x + 1 in
    if l then left x
    else right x;;

let move' l x = 
    if l then (fun y -> y - 1) x
    else (fun y -> y + 1) x

let rec map f l = match l with
    [] -> []
    | (h::t) -> (f h) :: (map f t);;

let rec fold f a l = match l with
    [] -> a
    | (h::t) -> fold f (f a h) t;;
