require 'socket'
s = TCPSocket.new 'localhost', 5555
s.write [0, 2, "sebastian", " \t ", "234",  " \t ", "\r\n"].pack("L<L<a*a3a*a3a2")

# Directives for pack
# Integer      | Array   |
# Directive    | Element | Meaning
# ---------------------------------------------------------------------------
#    C         | Integer | 8-bit unsigned (unsigned char)
#    S         | Integer | 16-bit unsigned, native endian (uint16_t)
#    L         | Integer | 32-bit unsigned, native endian (uint32_t)
#    Q         | Integer | 64-bit unsigned, native endian (uint64_t)
#    J         | Integer | pointer width unsigned, native endian (uintptr_t)
#              |         | (J is available since Ruby 2.3.)
#              |         |
#    c         | Integer | 8-bit signed (signed char)
#    s         | Integer | 16-bit signed, native endian (int16_t)
#    l         | Integer | 32-bit signed, native endian (int32_t)
#    q         | Integer | 64-bit signed, native endian (int64_t)
#    j         | Integer | pointer width signed, native endian (intptr_t)
#              |         | (j is available since Ruby 2.3.)
#              |         |
#    S_, S!    | Integer | unsigned short, native endian
#    I, I_, I! | Integer | unsigned int, native endian
#    L_, L!    | Integer | unsigned long, native endian
#    Q_, Q!    | Integer | unsigned long long, native endian (ArgumentError
#              |         | if the platform has no long long type.)
#              |         | (Q_ and Q! is available since Ruby 2.1.)
#    J!        | Integer | uintptr_t, native endian (same with J)
#              |         | (J! is available since Ruby 2.3.)
#              |         |
#    s_, s!    | Integer | signed short, native endian
#    i, i_, i! | Integer | signed int, native endian
#    l_, l!    | Integer | signed long, native endian
#    q_, q!    | Integer | signed long long, native endian (ArgumentError
#              |         | if the platform has no long long type.)
#              |         | (q_ and q! is available since Ruby 2.1.)
#    j!        | Integer | intptr_t, native endian (same with j)
#              |         | (j! is available since Ruby 2.3.)
#              |         |
#    S> L> Q>  | Integer | same as the directives without ">" except
#    J> s> l>  |         | big endian
#    q> j>     |         | (available since Ruby 1.9.3)
#    S!> I!>   |         | "S>" is same as "n"
#    L!> Q!>   |         | "L>" is same as "N"
#    J!> s!>   |         |
#    i!> l!>   |         |
#    q!> j!>   |         |
#              |         |
#    S< L< Q<  | Integer | same as the directives without "<" except
#    J< s< l<  |         | little endian
#    q< j<     |         | (available since Ruby 1.9.3)
#    S!< I!<   |         | "S<" is same as "v"
#    L!< Q!<   |         | "L<" is same as "V"
#    J!< s!<   |         |
#    i!< l!<   |         |
#    q!< j!<   |         |
#              |         |
#    n         | Integer | 16-bit unsigned, network (big-endian) byte order
#    N         | Integer | 32-bit unsigned, network (big-endian) byte order
#    v         | Integer | 16-bit unsigned, VAX (little-endian) byte order
#    V         | Integer | 32-bit unsigned, VAX (little-endian) byte order
#              |         |
#    U         | Integer | UTF-8 character
#    w         | Integer | BER-compressed integer

# Float        |         |
# Directive    |         | Meaning
# ---------------------------------------------------------------------------
#    D, d      | Float   | double-precision, native format
#    F, f      | Float   | single-precision, native format
#    E         | Float   | double-precision, little-endian byte order
#    e         | Float   | single-precision, little-endian byte order
#    G         | Float   | double-precision, network (big-endian) byte order
#    g         | Float   | single-precision, network (big-endian) byte order

# String       |         |
# Directive    |         | Meaning
# ---------------------------------------------------------------------------
#    A         | String  | arbitrary binary string (space padded, count is width)
#    a         | String  | arbitrary binary string (null padded, count is width)
#    Z         | String  | same as ``a'', except that null is added with *
#    B         | String  | bit string (MSB first)
#    b         | String  | bit string (LSB first)
#    H         | String  | hex string (high nibble first)
#    h         | String  | hex string (low nibble first)
#    u         | String  | UU-encoded string
#    M         | String  | quoted printable, MIME encoding (see RFC2045)
#    m         | String  | base64 encoded string (see RFC 2045, count is width)
#              |         | (if count is 0, no line feed are added, see RFC 4648)
#    P         | String  | pointer to a structure (fixed-length string)
#    p         | String  | pointer to a null-terminated string

# Misc.        |         |
# Directive    |         | Meaning
# ---------------------------------------------------------------------------
#    @         | ---     | moves to absolute position
#    X         | ---     | back up a byte
#    x         | ---     | null byte             
