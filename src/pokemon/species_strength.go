package pokemon

type SpeciesStrength [6]uint16

func (s SpeciesStrength) Sum() (sum uint32) {
	for _, v := range s {
		sum += uint32(v)
	}
	return sum
}

func (s SpeciesStrength) HP() uint16             { return s[0] }
func (s SpeciesStrength) Attack() uint16         { return s[1] }
func (s SpeciesStrength) Defense() uint16        { return s[2] }
func (s SpeciesStrength) SpecialAttack() uint16  { return s[3] }
func (s SpeciesStrength) SpecialDefense() uint16 { return s[4] }
func (s SpeciesStrength) Speed() uint16          { return s[5] }

func (s *SpeciesStrength) SetHP(v uint16)             { (*s)[0] = v }
func (s *SpeciesStrength) SetAttack(v uint16)         { (*s)[1] = v }
func (s *SpeciesStrength) SetDefense(v uint16)        { (*s)[2] = v }
func (s *SpeciesStrength) SetSpecialAttack(v uint16)  { (*s)[3] = v }
func (s *SpeciesStrength) SetSpecialDefense(v uint16) { (*s)[4] = v }
func (s *SpeciesStrength) SetSpeed(v uint16)          { (*s)[5] = v }
